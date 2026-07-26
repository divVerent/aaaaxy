$fn = 72;

player_offset = 2;
player_height = 22;
player_width = 10;
player_top_height = 3;
player_top_width1 = 6;
player_top_width2 = 2;
player_bottom_height = 2;
player_bottom_width1 = 8;
player_bottom_width2 = 6;
player_round_radius = 1;

foot_x = 3;
foot_y = -1.5;
foot_width = 4;
foot_length_x = 2;
foot_length_y = 3;

eye_x = 2.5;
eye_y = -4.5;
eye_z = 22;

eye_width = 4;
eye_height = 4;
eyebrow_x = -0.25;
eyebrow_y = 0.5;
eyebrow_z = 1;

pupil_width = 1;
pupil_height = 2;
pupil_x = 0.25;
pupil_y = -1.75;
pupil_z = 0;

bump_width = 2;
bump_height = 1;
bump_y = -5.75;
bump_z = 25.5;

module player_round() {
    minkowski() {
        children();
        sphere(r = player_round_radius);
    }
}

module player_base() {
    color("#0000aa")
    translate([0, 0, player_offset])
    player_round()
    cylinder(h = player_bottom_height, r1 = player_bottom_width2 / 2, r2 = player_bottom_width1 / 2);

    color("#5555ff")
    translate([0, 0, player_offset + player_bottom_height])
    player_round()
    cylinder(h = player_height, r = player_width / 2);

    color("#ffffff")
    translate([0, 0, player_offset + player_bottom_height + player_height])
    player_round()
    cylinder(h = player_top_height, r1 = player_top_width1 / 2, r2 = player_top_width2 / 2);
}

module player_foot(dir) {
    translate([dir * foot_x, foot_y, 0])
    hull() {
        translate([+1 * dir, -foot_length_y / 2, foot_width / 2])
        sphere(r = foot_width / 2);
        translate([+1 * dir, -foot_length_y / 2, 0])
        cylinder(r = foot_width / 2, h = foot_width / 2);
        translate([-1 * dir, +foot_length_y / 2, foot_width / 2])
        sphere(r = foot_width / 2);
        translate([-1 * dir, +foot_length_y / 2, 0])
        cylinder(r = foot_width / 2, h = foot_width / 2);
    }
}

module player_feet() {
    color("#aa5500")
    player_foot(-1);

    color("#ffff55")
    player_foot(+1);
}

module player_eye(dir) {
    translate([dir * eye_x, eye_y, eye_z]) {
        color("#555555")
        translate([dir * eyebrow_x, eyebrow_y, eyebrow_z])
        scale([eye_width, eye_width, eye_height] / 2)
        sphere(r = 1);

        color("#ffffff")
        scale([eye_width, eye_width, eye_height] / 2)
        sphere(r = 1);

        color("#000000")
        translate([dir * pupil_x, pupil_y, pupil_z])
        scale([pupil_width, pupil_width, pupil_height] / 2)
        sphere(r = 1);
    }
}

module player_eyes() {
    player_eye(+1);
    player_eye(-1);
}

module player_bump() {
    color("#55ffff")
    translate([0, bump_y, bump_z])
    hull() {
        translate([(bump_height - bump_width) / 2, 0, 0])
        sphere(r = bump_height / 2);
        translate([(bump_width - bump_height) / 2, 0, 0])
        sphere(r = bump_height / 2);
    }
}

module player() {
    player_base();
    player_feet();
    player_eyes();
    player_bump();
}

translate([0, 0, 0])
player();
