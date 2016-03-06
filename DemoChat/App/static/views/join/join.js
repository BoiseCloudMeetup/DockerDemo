(function() {
  "use strict";

  angular.module('app')
    .controller('joinController', controller);

  function controller($location, $routeParams, rooms) {
    var self = this,
    nickname = $routeParams.nickname;

    self.join = join;
    self.create = create;

    function join(room) {
      $location.path("/" + nickname + "/" + room);
    }

    function create() {
      rooms.random()
        .then(function(room) {
          console.log(room);
          $location.path("/" + nickname + "/" + room);
        }, function() {
          alert("Failed to create room ID");
        });
    }
  }

}());
