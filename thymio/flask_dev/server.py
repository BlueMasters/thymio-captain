from flask import Flask, jsonify, request
from flask_restful import Resource, Api
from threading import Thread
from time import sleep
from thymiolib import ThymioController
from datetime import datetime
import base64
import json
import shelve

app = Flask(__name__)
api = Api(app)

program = {}
running = False;
TIME_CM_CONV = 0.2
TIME_ANGLE_CONV = 0.025
WALL_VALUE = 3000

class progRunner(Thread):
    def __init__(self, prog):
        Thread.__init__(self)
        self.prog = prog
        self.askStop = False
        self.options = {
            "MoveForward": self.move_forward,
            "MoveBackward" : self.move_backward,
            "Turn" : self.turn,
            "FollowLine" : self.follow_line,
            "SetTopColor" : self.set_top_color,
            "SetBottomColor" : self.set_bottom_color
        }
        self.color = {
            "off": [0, 0, 0],
            "red": [32, 0, 0],
            "blue": [0, 0, 32],
            "green": [0, 32, 0],
            "pink": [32, 10, 10],
            "orange": [32, 20, 0],
            "white": [32, 32, 32]
        }
        self.thymio = ThymioController()

    def run(self):
        for instr in self.prog:
            action = instr["Action"]
            param = instr["Param"]
            self.options[action](param)
            if self.askStop:
                break
        global running
        self.thymio.set_motors(0,0)
        self.thymio.set_led_bottom_left(self.color["off"])
        self.thymio.set_led_bottom_right(self.color["off"])
        self.thymio.set_led_top(self.color["off"])
        running = False
        return

    def move_forward(self, param):
        self.thymio.set_motors(300,300)
        if param == "UntilWall":
            self.until_wall_front()
        elif param == "UntilBlackFloor":
            self.until_black_floor()
        elif param == "UntilWhiteFloor":
            self.until_white_floor()
        elif param == "10cm":
            self.until_time(10)
        elif param == "20cm":
            self.until_time(20)
        elif param == "50cm":
            self.until_time(50)
        return

    def move_backward(self, param):
        self.thymio.set_motors(-300,-300)
        if param == "UntilWall":
            self.until_wall_back()
        elif param == "UntilBlackFloor":
            self.until_black_floor()
        elif param == "UntilWhiteFloor":
            self.until_white_floor()
        else:
            self.until_time(int(param[:2]))
        return

    def turn(self, param):
        angle = 0
        if param[0] == 'R':
            angle = int(param[5:])
            self.thymio.set_motors(300, -300)
        else:
            angle = int(param[4:])
            self.thymio.set_motors(-300, 300)
        global TIME_ANGLE_CONV
        target_time = angle*TIME_ANGLE_CONV
        time = 0.0
        while self.askStop is False and time < target_time:
            sleep(0.1)
            time += 0.1
        return

    def follow_line(self, param):
        if param == "UntilWall":
            self.follow_line_until_wall()
        else:
            self.follow_line_until_time(int(param[:2]))

    def set_top_color(self, param):
        self.thymio.set_led_top(self.color[param])
        return

    def set_bottom_color(self, param):
        self.thymio.set_led_bottom_left(self.color[param])
        self.thymio.set_led_bottom_right(self.color[param])
        return

    def until_wall_front(self):
        while self.thymio.get_prox_h()[2] < WALL_VALUE and self.askStop is False:
            sleep(0.1)
        return

    def until_wall_back(self):
        while self.thymio.get_prox_h()[5] < 1000 and self.thymio.get_prox_h()[6] < 1000 and self.askStop == False:
            sleep(0.1)
        return

    def until_black_floor(self):
        while self.thymio.get_prox_v()[0] > 500 and self.thymio.get_prox_v()[1] > 500 and self.askStop is False:
            sleep(0.1)
        return

    def until_white_floor(self):
        while self.thymio.get_prox_v()[0] < 500 and self.thymio.get_prox_v()[1] < 500 and self.askStop is False:
            sleep(0.1)
        return

    def until_time(self,cm):
        global TIME_CM_CONV
        target_time = cm*TIME_CM_CONV
        time = 0.0
        while self.askStop is False and time < target_time:
            sleep(0.1)
            time += 0.1

    def follow_line_until_wall(self):
        left = True
        while self.thymio.get_prox_h()[2] < WALL_VALUE and self.askStop is False:
            ground = self.thymio.get_prox_v()
            if ground[0]<500 and ground[1]<500:
                self.thymio.set_motors(300,300)
            elif ground[0]<500 and ground[1]>500:
                self.thymio.set_motors(50,300)
                left = True
            elif ground[0]>500 and ground[1]<500:
                self.thymio.set_motors(300,50)
                left = False
            else:
                if left:
                    self.thymio.set_motors(-200,200)
                else:
                    self.thymio.set_motors(200,-200)
        self.thymio.set_motors(0,0)

    def follow_line_until_time(self,cm):
        left = True
        target_time = cm*TIME_CM_CONV
        time = 0.0
        while self.askStop is False and time < target_time:
            ground = self.thymio.get_prox_v()
            if ground[0]<500 and ground[1]<500:
                self.thymio.set_motors(300,300)
                time += 0.1
            elif ground[0]<500 and ground[1]>500:
                self.thymio.set_motors(50,300)
                left = True
                time += 0.06
            elif ground[0]>500 and ground[1]<500:
                self.thymio.set_motors(300,50)
                left = False
                time += 0.06
            else:
                if left:
                    self.thymio.set_motors(-200,200)
                else:
                    self.thymio.set_motors(200,-200)
        self.thymio.set_motors(0,0)


class calibLineRunner(Thread):
    def __init__(self):
        Thread.__init__(self)
        self.askStop = False
        self.thymio = ThymioController()

    def run(self):
        self.thymio.set_motors(300,300)
        while self.thymio.get_prox_v()[0] < 500 and self.thymio.get_prox_v()[1] < 500:
            sleep(0.1)
        a = datetime.now()
        while self.thymio.get_prox_v()[0] > 500 and self.thymio.get_prox_v()[1] > 500:
            sleep(0.1)
        b = datetime.now()
        self.thymio.set_motors(0,0)
        global TIME_CM_CONV
        TIME_CM_CONV = (b-a).total_seconds()/20.0
        global d
        d["cm"] = TIME_CM_CONV
        global running
        running = False


class calibRotRunner(Thread):
    def __init__(self):
        Thread.__init__(self)
        self.askStop = False
        self.thymio = ThymioController()

    def run(self):
        self.thymio.set_motors(300,-300)
        while self.thymio.get_prox_v()[0] < 500:
            sleep(0.1)
        a = datetime.now()
        while self.thymio.get_prox_v()[0] > 500:
            sleep(0.1)
        b = datetime.now()
        self.thymio.set_motors(0,0)
        global TIME_ANGLE_CONV
        TIME_ANGLE_CONV = (b-a).total_seconds()/90.0
        global d
        d["angle"] = TIME_ANGLE_CONV
        global running
        running = False

class calibWallRunner(Thread):
    def __init__(self):
        Thread.__init__(self)
        self.askStop = False
        self.thymio = ThymioController()

    def run(self):
        self.thymio.set_motors(100, 100)
        while self.thymio.get_prox_v()[0] < 500 and self.thymio.get_prox_v()[1] < 500:
            sleep(0.1)
        self.thymio.set_motors(0,0)
        global WALL_VALUE
        WALL_VALUE = self.thymio.get_prox_h()[2]
        global d
        d["wall"] = WALL_VALUE
        global running
        running = False

class upload(Resource):
    def put(self):
        global program
        encoded = request.get_json()["program"]
        decoded = base64.standard_b64decode(encoded)
        program = json.loads(decoded.decode("utf-8"))
        return {"result":"ok"}, 200


class run(Resource):
    def get(self):
        global program
        global crtThread
        global running
        if running:
            return {"result":"already running"},200
        else:
            running = True
            crtThread = progRunner(program)
            crtThread.start()
            return {"result":"ok"}, 200


class stop(Resource):
    def get(self):
        global crtThread
        crtThread.askStop = True
        # reinit tyhmio
        return {"result":"ok"}, 200

class prog(Resource):
    def get(self):
        global program
        return program,200

class state(Resource):
    def get(self):
        global running
        return {"running":running,"TIME_CM_CONV":TIME_CM_CONV,"TIME_ANGLE_CONV":TIME_ANGLE_CONV,"WALL_VALUE":WALL_VALUE},200

class ping(Resource):
    def get(self):
        return {"result":"ok"},200

class calibLine(Resource):
    def get(self):
        global running
        global crtThread
        if running:
            return {"result":"already running"},200
        else:
            running = True
            crtThread = calibLineRunner()
            crtThread.start()
            return {"result":"ok"},200

class calibRot(Resource):
    def get(self):
        global running
        global crtThread
        if running:
            return {"result":"already running"},200
        else:
            running = True
            crtThread = calibRotRunner()
            crtThread.start()
            return {"result":"ok"},200

class calibWall(Resource):
    def get(self):
        global running
        global crtThread
        if running:
            return {"result":"already running"},200
        else:
            running = True
            crtThread = calibWallRunner()
            crtThread.start()
            return {"result":"ok"},200


api.add_resource(upload, '/api/v1/upload')
api.add_resource(run, '/api/v1/run')
api.add_resource(stop, '/api/v1/stop')
api.add_resource(ping, '/api/v1/ping')
api.add_resource(prog, '/api/v1/prog')
api.add_resource(state, '/api/v1/state')
api.add_resource(calibLine, '/api/v1/calibline')
api.add_resource(calibRot, '/api/v1/calibrot')
api.add_resource(calibWall, '/api/v1/calibwall')

d = shelve.open("calib.conf")
if "cm" in d:
    TIME_CM_CONV = d["cm"]
if "angle" in d:
    TIME_ANGLE_CONV = d["angle"]
if "wall" in d:
    WALL_VALUE = d["wall"]

if __name__ == '__main__':
    app.run(debug=True)
