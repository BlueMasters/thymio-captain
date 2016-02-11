/*
 * @author   Lucy Linder <lucy.derlin@gmail.com>
 * @date     February 2016
 * @context  Thymio Captain
 *
 * Copyright 2016 BlueMagic. All rights reserved.
 * Use of this source code is governed by an Apache 2 license
 * that can be found in the LICENSE file.
 */
(function () {

    /**
     * @ngdoc service
     * @name thymioCaptain.rest.RestService
     * @description
     * Service to talk to the backend.
     */
    angular
        .module('thymioCaptain.rest')
        .factory('RestService', RestService);

    // --------------------------

    function RestService($resource, baseUrl) {

        return $resource('', {}, {


            /**
             * @ngdoc
             * @name getSession
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * get the status of a session
             *
             * @returns {httpPromise} resolves with the session, or fails with error description.
             */
            getSession: {method: 'GET', url: baseUrl + 'session/:id', params: {id: '@id'}},


            /**
             * @ngdoc
             * @name saveProgram
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * The body of the request is the JSON encoded program.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            saveProgram: {method: 'PUT', url: baseUrl + ':id/program', params: {id: '@id'}},

            /**
             * @ngdoc
             * @name getProgram
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * The body of the answer is the JSON encoded program.
             * @returns {httpPromise} resolves with the program, or fails with error description.
             */
            getProgram: {method: 'GET', url: baseUrl + ':id/program', params: {id: '@id'}},


            /**
             * @ngdoc
             * @name associateThymio
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Associates the robot with the session.
             * This command requires a valid "admin" cookie in the headers's request.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            associateThymio: {method: 'PUT', url: baseUrl + ':id/robot/:ip', params: {id: '@id', ip: '@ip'}},


            /**
             * @ngdoc
             * @name dissociateThymio
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Dissociates the robot with the session.
             * This command requires a valid "admin" cookie in the headers's request.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            dissociateThymio: {method: 'DEL', url: baseUrl + ':id/robot/:ip', params: {id: '@id', ip: '@ip'}},


            /**
             * @ngdoc
             * @name run
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Send the stopcommand to the associated robot. Returns an error if the robot is not associated.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            run: {method: 'GET', url: baseUrl + 'session/:id/stop', params: {id: '@id'}},

            /**
             * @ngdoc
             * @name uploadProgram
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Upload the code to the associated robot. Returns an error if the robot is not associated.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            upload: {method: 'GET', url: baseUrl + 'session/:id/upload', params: {id: '@id'}},


            /* ===================================================================*/
            /**
             * @ngdoc
             * @name isAdmin
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Checks if a connection is from an admin
             * @returns {httpPromise} resolves with a boolean or an error
             */
            isAdmin: {method: 'GET', url: baseUrl + 'isadmin'},

            /**
             * @ngdoc
             * @name getAssociations
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Returns a list of associations (admin only)
             * @returns {httpPromise} resolves with the list, or fails with error description.
             */
            upload: {method: 'GET', url: baseUrl + 'associations', params: {id: '@id'}, isArray: true}
        })
    }

}());