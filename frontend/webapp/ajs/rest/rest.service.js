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
     * @ngdoc service
     * @name thymioCaptain.rest.RestService
     * @description
     * Service to talk to the backend.
     */
    angular
        .module( 'thymioCaptain.rest' )
        .factory( 'RestService', RestService );

    // --------------------------

    function RestService( $resource, baseUrl ){

        return $resource( '', {}, {


            /**
             * @ngdoc
             * @name infos
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Returns information based on the cookie.
             *
             * @returns {httpPromise} resolves with the object {cardId: "", isadmin: bool}, or fails with error
             * description.
             */
            infos: {method: 'GET', url: baseUrl + 'info'},


            /**
             * @ngdoc
             * @name getCardData
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Returns the card's data based on this cardId.
             * @returns {httpPromise} resolves with the object {cardId: "", program: ""}, or fails with error
             * description.
             */
            getCardData: {method: 'GET', url: baseUrl + 'card/:cardId', params: {cardId: '@cardId'}},

            /**
             * @ngdoc
             * @name saveProgram
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * The body of the request is the object {cardId: "", program: ""}
             * @returns {httpPromise} resolves, or fails with error description.
             */
            setCardData: {method: 'PUT', url: baseUrl + 'card/:cardId', params: {cardId: '@cardId'}},


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
            associateThymio: {
                method: 'PUT',
                url   : baseUrl + 'robot/:name/card/:cardId',
                params: {cardId: '@cardId', name: '@name'}
            },


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
            dissociateThymio: {
                method: 'DELETE',
                url   : baseUrl + 'robot/:name/card/:cardId',
                params: {cardId: '@cardId', name: '@name'}
            },

            /**
             * @ngdoc
             * @name run
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Send the run command to the associated robot. Returns an error if the robot is not associated.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            run: {method: 'GET', url: baseUrl + 'card/:cardId/run', params: {cardId: '@cardId'}},

            /**
             * @ngdoc
             * @name stop
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Send the stop command to the associated robot. Returns an error if the robot is not associated.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            stop: {method: 'GET', url: baseUrl + 'card/:cardId/stop', params: {cardId: '@cardId'}},

            /**
             * @ngdoc
             * @name upload
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Upload the code to the associated robot. Returns an error if the robot is not associated.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            upload: {method: 'GET', url: baseUrl + 'card/:cardId/upload', params: {cardId: '@cardId'}},



            /* ===================================================================*/

            /**
             * @ngdoc
             * @name robotStatus
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Returns the status of a robot.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            robotStatus: {method: 'GET', url: baseUrl + 'robot/:name', params: {name: '@name'}},

            /**
             * @ngdoc
             * @name addRobot
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Add or update a robot. Payload: {"url": "<robot url>"}
             * @returns {httpPromise} resolves, or fails with error description.
             */
            addRobot: {method: 'PUT', url: baseUrl + 'robot/:name', params: {name: '@name'}},

            /**
             * @ngdoc
             * @name robotStatus
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Delete the robot
             * @returns {httpPromise} resolves, or fails with error description.
             */
            deleteRobot: {method: 'DEL', url: baseUrl + 'robot/:name', params: {name: '@name'}},


            /**
             * @ngdoc
             * @name getRobots
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * Returns the list of all known robots with the associated card (admin only)
             * @returns {httpPromise} resolves with the list (url,name,cardId), or fails with error description.
             */
            getRobots: {method: 'GET', url: baseUrl + 'robots', isArray: true},

            /**
             * @ngdoc
             * @name pingRobot
             * @methodOf thymioCaptain.rest.RestService
             *
             * @description
             * check if the robot is connected.
             * @returns {httpPromise} resolves, or fails with error description.
             */
            pingRobot: {method: 'GET', url: baseUrl + 'robot/:name/ping', params: {name: '@name'}}
        } )
    }

}());