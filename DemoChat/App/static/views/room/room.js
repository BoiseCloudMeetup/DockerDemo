(function() {
  "use strict";

  angular.module('app')
    .controller('roomController', controller);

  function controller($scope, $routeParams, $http, $location, $interval, users, messages) {

    var self = this;

    var nickname = $routeParams.nickname,
        roomID = $routeParams.roomID;

    // Open web socket connection
    var host = window.location.host;
    console.log(host);
    var socket = new WebSocket("ws://" + host + "/ws/?roomID=" + roomID + "&nickname=" + nickname);
    socket.onopen = function() {
      console.log("socket opened");
    };
    socket.onmessage = function(event) {
      var msg = JSON.parse(event.data);
      $scope.$apply(function() {
        switch(msg.Action) {
          case "userenter":
            self.users.push(msg.Content); break;
          case "userexit":
            users.get(roomID)
              .then(function(response) {
                self.users = response;
              }, function() {
                console.log("failed to get users");
              });
            break;
          case "message":
            self.messages.push(msg.Content); break;
        }
      });
      console.log(event.data);
    };

    // Handle messaging
    self.messages = [];
    self.users = [];
    self.text = '';
    self.sendMessage = sendMessage;

    messages.get(roomID)
      .then(function(response) {
        self.messages = response;
      });

    var userPoller = $interval(function() {
      users.get(roomID)
      .then(function(response) {
        self.users = response;
      }, function() {
        console.log("failed to get users");
      });
    }, 1000);

    function sendMessage(text) {
      if (text === "") {
        return;
      }

      var message = {RoomID: roomID, Nickname: nickname, Text: text};

      $http.post("/ws/send", message)
        .then(function() {
          console.log("ws message sent");
        }, function() {
          console.log("ws message failed to send");
        });

      self.text = "";
      window.setTimeout(function () {
        updateScroll();
      }, 1);
    }


    function updateScroll(){
      var element = document.getElementById("scroller");
      element.scrollTop = element.scrollHeight;
    }

    $scope.$on("$destroy", function() {
      socket.close();
      $interval.cancel(userPoller);
    });

  }

}());
