(function() {
  "use strict";

  angular.module("app", ["ngRoute"])
    .config(function($routeProvider) {
      $routeProvider
        .when("/", {
          templateUrl: "/views/login/login.html",
          controller: "loginController",
          controllerAs: "ctl",
        })
        .when("/:nickname", {
          templateUrl: "/views/join/join.html",
          controller: "joinController",
          controllerAs: "ctl",
        })
        .when("/:nickname/:roomID", {
          templateUrl: "/views/room/room.html",
          controller: "roomController",
          controllerAs: "ctl",
        })
        .otherwise({
          redirectTo: "/"
        });
    });



}());
