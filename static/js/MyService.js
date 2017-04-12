var app = angular.module('myApp')

app.service('RoleSrv', function() {
	var role = "admin";
    this.setRole = function (r) {
        role=r;
    }
	this.getRole=function(){
		return role;
	}
	
});