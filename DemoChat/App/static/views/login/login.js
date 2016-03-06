(function() {
  "use strict";

  angular.module('app')
    .controller('loginController', controller);

    function controller($location) {
      var self = this;

      self.login = login;

      function login(name) {
        $location.path("/" + name);
      }

    }

}());
