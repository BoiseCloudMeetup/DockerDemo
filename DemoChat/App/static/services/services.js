(function() {
  "use strict";

  angular.module("app")
    .factory("rooms", roomService)
    .factory("messages", messageService)
    .factory("users", userService);


  function roomService($http, $q) {
    var service = {
      random: function() {
        return $http.get("/api/keys/")
          .then(function(response) {
            return response.data;
          }, function() {
            return $q.reject();
          });
      }
    };
    return service;
  }

  function messageService($http) {
    var service = {
      get: function(roomID) {
        return $http.get("/api/messages/?roomID=" + roomID)
          .then(function(response) {
            return response.data;
          }, function(){
            console.log("Couldn't get messages.");
          });
      }
    };
    return service;
  }

  function userService($http) {
    var service = {
      get: function(roomID) {
        return $http.get("/api/users/?roomID=" + roomID)
          .then(function(response) {
            return response.data;
          }, function() {
            console.log("Couldn't get users.");
          });
      }
    };
    return service;
  }

}());
