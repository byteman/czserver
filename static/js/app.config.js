'use strict';

angular.
	module('myApp').
	config(function ($routeProvider) {
    $routeProvider.
      when('/weight/:id?', {
        templateUrl: 'weight.html',
        controller: 'weightCtrl'
    }).
    when('/online', {
        templateUrl: 'online.html',
        controller: 'OnLineCtrl'
    }).
	when('/param', {
        templateUrl: 'param.html',
        controller: 'ParamCtrl'
    }).
   when('/gpslist', {
        templateUrl: 'gpslist.html',
        controller: 'GpsListCtrl'
    }). 
	when('/update', {
        templateUrl: 'update.html',
        controller: 'UpdateCtrl'
    }).
	when('/gps', {
        templateUrl: 'gps2.html',
        controller: 'GpsCtrl'
    }).
    otherwise({
        redirectTo: '/online'
    });
	

});

