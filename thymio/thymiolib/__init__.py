import httplib2
import json

"""
@author: Damien Goetschi
@organiation: Haute ecole d'ingenierie et d'architecture Fribourg

"""

class ThymioController:
    """
    This class allows the user to control all leds, motors and sensors of a thymio

    """
    __errorConnReset = "Error: Connection reset by peer. Please check the Thymio is well connected to the Raspberry Pi"
    __errorConnRefused = "Error: Connection refused. Please check the Thymio is well connected to the Raspberry Pi or " \
                         "you entered the correct IP address"
    __errorNotFound = "Server not found. Please check you entered the correct IP address"
    __errorNoRoute = "No route to host. Please check you entered the correct IP address"

    __ip = "localhost"
    """IP address on which asebahttp runs"""

    __port = 3000
    """port on which asebahttp runs"""

    __rootUrl = ""
    """URL build with ip and port"""

    __debug = False
    """Debug mode"""

    __http = httplib2.Http()
    """Http objet to call RESTful API"""

    def __init__(self, ip="localhost", port=3000, debug=False):
        """
        Create a thymio controller with an ip address and a given port

        @param ip: asebahttp's ip
        @type ip:  string
        @param port: asebahttp's port
        @type port: int
        @param debug: true if debug mode is active (print detail about errors)
        @type debug: bool
        """
        self.__ip = ip
        self.__port = port
        self.__debug = debug
        self.__rootUrl = "http://{0}:{1}/nodes/thymio-II/".format(self.__ip, self.__port)

    def get_motors(self):
        """Get real speed for each motor

        This returns the speed of both thymio motors.

        @return: A list of two values containing the left and right motor values
        @rtype: list of int [left, right] (empty list in case of error)

        @note: \n

        >>> t.get_motors()
        [-52, 137]
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"motor.left.speed", "GET")
            left = (json.loads(content.decode('utf-8'))[0])
            response, content = self.__http.request(self.__rootUrl+"motor.right.speed", "GET")
            right = (json.loads(content.decode('utf-8'))[0])
            return [left, right]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return []
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return []
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return []
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return []

    def set_motors(self, left, right):
        """Change target speed for each motor

        This changes the speed of both motors  and return True on success.

        @param left: Left motor speed
        @param right: Right motor speed
        @type left: int (between -500 and 500)
        @type right: int (between -500 and 500)
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> t.set_motors(500, -500)
        True
        >>> t.set_motors(-1000, 0)
        False
        >>> t.set_motors(0, 0)
        True
        """
        if type(left)!=int or type(right)!=int:
            if self.__debug:
                print("Error: wrong type\n"+self.set_motors.__doc__)
            return False
        if left > 500 or left < -500 or right > 500 or right < -500:
            if self.__debug:
                print("Error: params not in range\n"+self.set_motors.__doc__)
            return False
        try:
            self.__http.request(self.__rootUrl+"eventMotors/{0}/{1}".format(left,right), "POST")
            return True
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def get_prox_h(self):
        """Get values for horizontal ir sensor.

        This returns the values of each horizontal ir sensor.

        @return: A value between 0 and 5000 for each of 7 sensors.
        The higher the value, the closer an object

        [0-4]: front sensors left to right
        [5-6]: back sensors left to right
        @rtype: list of int (empty list in case of error)

        @note: \n

        >>> t.get_prox_h()
        [0, 2701, 2579, 2489, 0, 3632, 4490]
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"prox.horizontal")
            return json.loads(content.decode('utf-8'))
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return []
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return []
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return []
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return []

    def get_prox_v(self):
        """Get values for vertical (ground) ir sensor.

        This returns the values of each vertical ir sensor.

        @return: A value between 0 and 1000 for each of 2 sensors.
        The higher the value, the lighter the ground.

        [0-1]: ground sensors left to right
        @rtype: list of int (empty list in case of error)

        @note: \n

        >>> t.get_prox_v()
        [980,240]
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"prox.ground.delta", "GET")
            return json.loads(content.decode('utf-8'))
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return []
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return []
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return []
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return []

    def get_acc(self):
        """Get values of accelerometer.

        This returns the values of each axis.

        @return: A value for each of 3 axis.

        [roll,pitch,yaw]
        @rtype: list of int (empty list in case of error)

        @note: \n

        >>> t.get_acc()
        [1,-1,23]
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"acc")
            return json.loads(content.decode('utf-8'))
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return []
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return []
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return []
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return []

    def get_button_backward(self):
        """Get value of backward button.

        This returns the current state of the button

        @return: 1 if button is pressed, 0 otherwise; -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_button_backward()
        0
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"button.backward")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1

    def get_button_center(self):
        """Get value of center button.

        This returns the current state of the button

        @@return: 1 if button is pressed, 0 otherwise; -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_button_center()
        1
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"button.center")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1

    def get_button_forward(self):
        """Get value of forward button.

        This returns the current state of the button

        @@return: 1 if button is pressed, 0 otherwise; -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_button_forward()
        0
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"button.forward")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1

    def get_button_left(self):
        """Get value of left button.

        This returns the current state of the button

        @@return: 1 if button is pressed, 0 otherwise; -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_button_left()
        1
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"button.left")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1

    def get_button_right(self):
        """Get value of right button.

        This returns the current state of the button

        @@return: 1 if button is pressed, 0 otherwise; -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_button_right()
        0
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"button.right")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1


    def get_mic_intensity(self):
        """Get value of mic.

        This returns the intensity of mic.

        @return: Between 0 and 255: the intensit;. -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_mic_intensity()
        129
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"mic.intensity")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1

    def get_temperature(self):
        """Get value of temperature.

        This returns the temperature of sensor

        @return: Temperature in tenths of a degree Celsius; -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_temperature()
        312
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"temperature")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1

    def get_rc_last_command(self):
        """Get last rc command received.

        @return: Command number (between 0 and 127); -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_rc_last_command()
        80
        """
        try:
            self.__http.request(self.__rootUrl+"new_rc/0","POST")
            response, content = self.__http.request(self.__rootUrl+"rc5.command")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1

    def get_rc_last_address(self):
        """Get last rc address received.

        @return: Address number (between 0 and 31). -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_rc_last_address()
        0
        """
        try:
            self.__http.request(self.__rootUrl+"new_rc/0","POST")
            response, content = self.__http.request(self.__rootUrl+"rc5.address")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1

    def get_rc_new(self):
        """Get if new rc code received.

        @return: 1 if new rc code received, 0 otherwise; -1 in case of error

        @rtype: int

        @note: \n

        >>> t.get_temperature()
        129
        """
        try:
            response, content = self.__http.request(self.__rootUrl+"new_rc")
            return json.loads(content.decode('utf-8'))[0]
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return -1
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return -1
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return -1
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return -1

    def set_sound_system(self,sound):
        """Set a sound system

        @param sound: sound to play:
        -1: stop playing sound
        0: startup sound
        1: shutdown sound
        2: arrow button sound
        3: central button sound
        4: free-fall (scary) sound
        5: collision sound
        6: target ok for friendly behaviour
        7: target detect for friendly behaviour
        @type sound: int
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_sound_system(-1)
        True
        >>> set_sound_system(20)
        False
        """
        if type(sound)!=int:
            if self.__debug:
                print("Error: wrong type\n"+self.set_sound_system.__doc__)
            return False
        if sound > 7 or sound < -1:
            if self.__debug:
                print("Error: params not in range\n"+self.set_sound_system.__doc__)
            return False
        try:
            response, content = self.__http.request(self.__rootUrl+"eventSoundSystem/{0}".format(sound))
            return True
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_sound_freq(self,freq,ds):
        """Set a sound frequency

        @param freq: sound frequency (Hz)
        @type freq: float (between 0 and 7812.5)
        @param ds:  sound duration in 1/60s. Specifying a 0 duration plays the sound continuously and specifying a -1
        duration stops the sound.
        @type ds: float
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_sound_freq(200, 60)
        True
        """
        if type(freq)!=float or type(ds)!=float:
            if self.__debug:
                print("Error: wrong type\n"+self.set_sound_freq.__doc__)
            return False
        if freq > 7812.5 or freq < 0 or ds < -1:
            if self.__debug:
                print("Error: params not in range\n"+self.set_sound_freq.__doc__)
            return False
        try:
            response, content = self.__http.request(self.__rootUrl+"eventSoundFreq/{0}/{1}".format(freq,ds))
            return True
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_led_top(self,r,g,b):
        """Set color for led on top

        @param r: value of red (0-32)
        @type r: int
        @param g: value of green (0-32)
        @type g: int
        @param b: value of blue (0-32)
        @type b: int
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_led_top(10,0,32)
        True
        """
        if type(r)!=int or type(g)!=int or type(b)!=int:
            if self.__debug:
                print("Error: wrong type\n"+self.set_led_top.__doc__)
            return False
        if r > 32 or r < 0 or g > 32 or g < 0 or b > 32 or b < 0:
            if self.__debug:
                print("Error: params not in range\n"+self.set_led_top.__doc__)
            return False
        try:
            response, content = self.__http.request(self.__rootUrl+"eventLedTop/{0}/{1}/{2}".format(r,g,b))
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_led_bottom_left(self,color):
        """Set color for led on bottom left

        @param r: value of red (0-32)
        @type r: int
        @param g: value of green (0-32)
        @type g: int
        @param b: value of blue (0-32)
        @type b: int
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_led_bottom_left(10,0,32)
        True
        """
        r = color[0]
        g = color[1]
        b = color[2]
        if type(r)!=int or type(g)!=int or type(b)!=int:
            if self.__debug:
                print("Error: wrong type\n"+self.set_led_bottom_left.__doc__)
            return False
        if r > 32 or r < 0 or g > 32 or g < 0 or b > 32 or b < 0:
            if self.__debug:
                print("Error: params not in range\n"+self.set_led_bottom_left.__doc__)
            return False
        try:
            response, content = self.__http.request(self.__rootUrl+"eventLedBotLeft/{0}/{1}/{2}".format(r,g,b))
            return True
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_led_bottom_right(self, color):
        """Set color for led on bottom right

        @param r: value of red (0-32)
        @type r: int
        @param g: value of green (0-32)
        @type g: int
        @param b: value of blue (0-32)
        @type b: int
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_led_bottom_right(10,0,32)
        True
        """
        r = color[0]
        g = color[1]
        b = color[2]
        if type(r)!=int or type(g)!=int or type(b)!=int:
            if self.__debug:
                print("Error: wrong type\n"+self.set_led_bottom_right.__doc__)
            return False
        if r > 32 or r < 0 or g > 32 or g < 0 or b > 32 or b < 0:
            if self.__debug:
                print("Error: params not in range\n"+self.set_led_bottom_right.__doc__)
            return False
        try:
            response, content = self.__http.request(self.__rootUrl+"eventLedBotRight/{0}/{1}/{2}".format(r,g,b))
            return True
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_led_temp(self, color):
        """Set color for temperature led on right

        @param r: value of red (0-32)
        @type r: int
        @param b: value of blue (0-32)
        @type b: int
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_led_temp(10,32)
        True
        """
        r = color[0]
        b = color[1]
        if type(r)!=int or type(b)!=int:
            if self.__debug:
                print("Error: wrong type\n"+self.set_led_temp.__doc__)
            return False
        if r > 32 or r < 0 or b > 32 or b < 0:
            if self.__debug:
                print("Error: params not in range\n"+self.set_led_temp.__doc__)
            return False

        try:
            response, content = self.__http.request(self.__rootUrl+"eventLedTemp/{0}/{1}".format(r, b))
            return True
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_led_sound(self, r):
        """Set color for sound led on left

        @param r: value of red (0-32)
        @type r: int
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_led_sound(20)
        True
        """
        if type(r)!=int:
            if self.__debug:
                print("Error: wrong type\n"+self.set_led_sound.__doc__)
            return False
        if r > 32 or r < 0:
            if self.__debug:
                print("Error: params not in range\n"+self.set_led_sound.__doc__)
            return False
        try:
            response, content = self.__http.request(self.__rootUrl+"eventLedSound/{0}".format(r))
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_led_button(self, leds):
        """Set color for led on top

        @param leds: value for leds for each arrow button (forward, right, backward, left)
        @type leds: int[4]
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_led_button([0, 32, 0, 32])
        True
        """
        if len(leds)!=4:
            if self.__debug:
                print("Error: incorrect number of values in array\n"+self.set_led_button.__doc__)
            return False
        for led in leds:
            if type(led)!=int:
                if self.__debug:
                    print("Error: wrong type\n"+self.set_led_button.__doc__)
                return False
            if led > 32 or led < 0:
                if self.__debug:
                    print("Error: params not in range\n"+self.set_led_button.__doc__)
                return False
        try:
            response, content = self.__http.request(self.__rootUrl+"eventLedButton/{0}/{1}/{2}/{3}".format(leds[0], leds[1],
                                                                                                   leds[2], leds[3]))
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_led_circle(self, leds):
        """Set color for led circle on top

        @param leds: value for leds for each part of circle (clockwise from forward)
        @type leds: int[8]
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_led_circle([0, 32, 0, 32, 0, 32, 0, 32])
        True
        """

        if len(leds)!=8:
            if self.__debug:
                print("Error: incorrect number of values in array\n"+self.set_led_circle.__doc__)
            return False
        for led in leds:
            if type(led)!=int:
                if self.__debug:
                    print("Error: wrong type\n"+self.set_led_circle.__doc__)
                return False
            if led > 32 or led < 0:
                if self.__debug:
                    print("Error: params not in range\n"+self.set_led_circle.__doc__)
                return False
        try:
            self.__http.request(self.__rootUrl+"eventLedCircle/{0}/{1}/{2}/{3}/{4}/{5}/{6}/{7}".
                                        format(leds[0], leds[1], leds[2], leds[3], leds[4], leds[5], leds[6], leds[7]))
            return True
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_led_prox_h(self, leds):
        """Set color for leds of horizontal proximity sensors

        @param leds: value for leds for each proximity sensor (0-5: front left to right, 6-7 back left to right)
        @type leds: int[8]
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_led_prox_h([0, 0, 32, 32, 0, 0, 16, 16])
        True
        """
        if len(leds)!=8:
            if self.__debug:
                print("Error: incorrect number of values in array\n"+self.set_led_prox_h.__doc__)
            return False
        for led in leds:
            if type(led)!=int:
                if self.__debug:
                    print("Error: wrong type\n"+self.set_led_prox_h.__doc__)
                return False
            if led > 32 or led < 0:
                if self.__debug:
                    print("Error: params not in range\n"+self.set_led_prox_h.__doc__)
                return False
        try:
            self.__http.request(self.__rootUrl+"eventLedProxH/{0}/{1}/{2}/{3}/{4}/{5}/{6}/{7}".
                                        format(leds[0], leds[1], leds[2], leds[3], leds[4], leds[5], leds[6], leds[7]))
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def set_led_prox_v(self, leds):
        """Set color for leds of vertical proximity sensors

        @param leds: value for leds for each proximity sensor (0: left, 1: right)
        @type leds: int[2]
        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> set_led_prox_v([32, 32])
        True
        """
        if len(leds)!=2:
            if self.__debug:
                print("Error: incorrect number of values in array\n"+self.set_led_prox_v.__doc__)
            return False
        for led in leds:
            if type(led)!=int:
                if self.__debug:
                    print("Error: wrong type\n"+self.set_led_prox_v.__doc__)
                return False
            if led > 32 or led < 0:
                if self.__debug:
                    print("Error: params not in range\n"+self.set_led_prox_v.__doc__)
                return False
        try:
            response, content = self.__http.request(self.__rootUrl+"eventLedProxV/{0}/{1}".format(leds[0], leds[1]))
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False

    def reset(self):
        """Set motors to 0 and turn off top and bottom leds

        @return: True if the method terminated correctly
        @rtype: bool

        @note: \n

        >>> reset()
        True
        """
        try:
            self.__http.request(self.__rootUrl+"eventReset")
        except httplib2.ServerNotFoundError:
            if self.__debug:
                print(self.__errorNotFound)
            return False
        except ConnectionRefusedError:
            if self.__debug:
                print(self.__errorConnRefused)
            return False
        except ConnectionResetError:
            if self.__debug:
                print(self.__errorConnReset)
            return False
        except OSError:
            if self.__debug:
                print(self.__errorNoRoute)
            return False
