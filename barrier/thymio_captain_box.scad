$fn = 50;
w   =  4;   // wood width
adv = 80;   // box extension
x0  = 10.5; // first hole/notch
dx  = 20;   // holes/notch spacing

nh = 6;     // number of holes
hx = 8;     // hole width
lw = 0.01;  // laser width

bar_len    = 150;
bar_width  =  40;
bar_height =  60;

box_width  = 100;
box_height =  90;
box_r      =  15;

module box_base(width, height, radius) {
    union() {
        square(size = [width, height-radius]);
        translate([radius,height-radius]){
            circle(r=radius);
        }
        translate([width-radius,height-radius]){
            circle(r=radius);
        }
        translate([radius,height-2*radius,0]){
            square(size = [width-2*radius, 2*radius]);
        }
    }
}

module top(bx) {
    difference() {
        union() {
            square(size = [bar_len, bar_width]);
            if (bx) {
                translate([(bar_len-box_width)/2,0,0]){
                    box_base(box_width, box_height, box_r);
                }
            }         
        }
        union() {
            for (i = [0:nh]) {
                translate([x0+i*dx+lw,-lw,0]) {
                    square(size = [hx-lw, w+2*lw]);
                }
                translate([x0+i*dx+lw,bar_width-w+lw,0]) {
                    square(size = [hx-lw, w+2*lw]);
                }
            }
        }
    }
}

module side() {
    union() {
        square(size = [bar_len, bar_height]);
        for (i = [0:nh]) {
            translate([x0+i*dx,-w-2*lw,0]) {
                square(size = [hx+lw, w+2*lw]);
            }
            translate([x0+i*dx,bar_height-2*lw,0]) {
                square(size = [hx+lw, w+2*lw]);
            }
        }
    }
}

top(true);
translate([0, -41, 0]) {
    top(false);
}
translate([0, -106, 0]) {
    side();
}
translate([0, -175, 0]) {
    side();
}
