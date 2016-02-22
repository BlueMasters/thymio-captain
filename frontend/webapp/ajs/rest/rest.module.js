/*
 * @author   Lucy Linder <lucy.derlin@gmail.com>
 * @date     February 2016
 * @context  Thymio Captain
 *
 * Copyright 2016 BlueMagic. All rights reserved.
 * Use of this source code is governed by an Apache 2 license
 * that can be found in the LICENSE file.
 */
(function(){

    /**
     * @ngdoc overview
     * @name thymioCaptain.rest
     * @requires $resource
     * @description
     * This module handles the interaction with the server and the app front end.
     */
    angular
        .module( 'thymioCaptain.rest', ['ngResource', 'base64', 'ngCookies'] )
        .config( function( $httpProvider, $base64 ){

            // handle authentication cookie
            // note: service can't be injected directly into config, so I do it manually
            // more info at: http://stackoverflow.com/questions/15358029/why-am-i-unable-to-inject-angular-cookies
            angular.injector( ['ngCookies'] ).invoke( ['$cookies', function( $cookies ){
                var auth = $cookies.get( "session-key" );
                if( auth ){
                    $httpProvider.defaults.header['Authorization'] = "Cookie " + auth;
                }
            }] );


            // handle the from/to base64 (program argument only)
            $httpProvider.defaults.transformRequest.unshift( function( data, headerGetter ){
                if( data && data.program ){
                    data.program = $base64.encode( JSON.stringify( data.program ) );
                }
                return data;
            } );

            $httpProvider.defaults.transformResponse.push( function( data, headerGetter ){
                if( data && data.program ){
                    data.program = JSON.parse( $base64.decode( data.program ) );
                }
                return data;
            } );

        } );


}());
