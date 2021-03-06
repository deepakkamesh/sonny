// Constants for roomba.
var CHARGING_STATE = {
    0: "No Charge",
    1: "Recond.",
    2: "Charging",
    3: "Trickle",
    4: "Waiting",
    5: "Fault",
}

var OI_MODE = {
    0: "Off",
    1: "Passive",
    2: "Safe",
    3: "Full",
}

var IR_CODE_NAMES = {
    0: "No IR Detected",
    161: "Force Field",
    164: "Green Buoy",
    165: "Green Buoy and Force Field",
    168: "Red Buoy",
    169: "Red Buoy and Force Field",
    172: "Red and Green Buoy",
    173: "Red, Green Buoy and Force Field",
}

$(document).ready(function() {
    var count = 0;
    var motorBackButton = document.querySelector('#motor-back');
    motorBackButton.addEventListener('click', function() {
        $.post('/api/move/_?dir=bwd', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });

    $(document).keydown(function(e) {
        if (e.which === 40) {
            $.post('/api/move/_?dir=bwd', "", function(data, status) {
                if (data.Err != '') {
                    console.log(data.Err);
                    return
                }
            });
        }
    });

    var motorForwardButton = document.querySelector('#motor-forward');
    motorForwardButton.addEventListener('click', function() {
        $.post('/api/move/_?dir=fwd', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });

    $(document).keydown(function(e) {
        if (e.which === 38) {
            $.post('/api/move/_?dir=fwd', "", function(data, status) {
                if (data.Err != '') {
                    console.log(data.Err);
                    return
                }
            });
        }
    });

    var motorLeftButton = document.querySelector('#motor-left');
    motorLeftButton.addEventListener('click', function() {
        $.post('/api/move/_?dir=left', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });

    $(document).keydown(function(e) {
        if (e.which === 37) {
            $.post('/api/move/_?dir=left', "", function(data, status) {
                if (data.Err != '') {
                    console.log(data.Err);
                    return
                }
            });
        }
    });

    var motorRightButton = document.querySelector('#motor-right');
    motorRightButton.addEventListener('click', function() {
        $.post('/api/move/_?dir=right', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });

    $(document).keydown(function(e) {
        if (e.which === 39) {
            $.post('/api/move/_?dir=right', "", function(data, status) {
                if (data.Err != '') {
                    console.log(data.Err);
                    return
                }
            });
        }
    });



    var driveVelSel = document.querySelector('#drive_velocity_sel');
    driveVelSel.addEventListener('click', function() {
        val = $('#drive_velocity_sel').val();
        $("#drive_velocity_sel_disp").empty()
        $("#drive_velocity_sel_disp").append(val);

        $.post('/api/setparam/_?velocity=' + val, "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });

    });


    // Servo Controls.
    var servoDownButton = document.querySelector('#servo-down');
    servoDownButton.addEventListener('click', function() {
        $.post('/api/servorotate/_?dir=down', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
            $("#servo-up-tip").empty();
            $("#servo-up-tip").append(data.Data['vert']);
            $("#servo-down-tip").empty();
            $("#servo-down-tip").append(data.Data['vert']);
        });
    });

    var servoUpButton = document.querySelector('#servo-up');
    servoUpButton.addEventListener('click', function() {
        $.post('/api/servorotate/_?dir=up', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
            $("#servo-down-tip").empty();
            $("#servo-down-tip").append(data.Data['vert']);
            $("#servo-up-tip").empty();
            $("#servo-up-tip").append(data.Data['vert']);
        });
    });

    var servoLeftButton = document.querySelector('#servo-left');
    servoLeftButton.addEventListener('click', function() {
        $.post('/api/servorotate/_?dir=left', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
            $("#servo-left-tip").empty();
            $("#servo-left-tip").append(data.Data['horiz']);
            $("#servo-right-tip").empty();
            $("#servo-right-tip").append(data.Data['horiz']);
        });
    });

    var servoRightButton = document.querySelector('#servo-right');
    servoRightButton.addEventListener('click', function() {
        $.post('/api/servorotate/_?dir=right', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
            $("#servo-left-tip").empty();
            $("#servo-left-tip").append(data.Data['horiz']);
            $("#servo-right-tip").empty();
            $("#servo-right-tip").append(data.Data['horiz']);
        });
    });

    var servoAngleDeltaSel = document.querySelector('#servo_angle_step');
    servoAngleDeltaSel.addEventListener('click', function() {
        val = $('#servo_angle_step').val();
        $("#servo_angle_step_disp").empty();
        $("#servo_angle_step_disp").append(val);

        $.post('/api/setparam/_?servoDelta=' + val, "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });

    var i2cEnBtn = document.querySelector('#i2c_en');
    i2cEnBtn.addEventListener('click', function() {
        var param = '';
        if (document.getElementById('i2c_en').checked) {
            param = 'on';
        } else {
            param = 'off';
        }

        $.post('/api/i2c_en/_?param=' + param, "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });

    var auxPwrBtn = document.querySelector('#aux_power');
    auxPwrBtn.addEventListener('click', function() {
        var param = '';
        if (document.getElementById('aux_power').checked) {
            param = 'on';
        } else {
            param = 'off';
        }

        $.post('/api/roomba_cmd/_?cmd=aux_power&param=' + param, "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });

    /*
        var roombaPowerBtn = document.querySelector('#roomba_power');
        roombaPowerBtn.addEventListener('click', function() {
            var action = '';
            if (document.getElementById('roomba_power').checked) {
                action = 'power_on';
            } else {
                action = 'power_off';
            }

            $.post('/api/roomba_cmd/_?cmd=' + action, "", function(data, status) {
                if (data.Err != '') {
                    console.log(data.Err);
                    return
                }
            });
        });*/

    $("#mode_full, #mode_safe, #mode_passive").change(function() {
        mode = $('input[name=roomba_mode]:checked').val();
        $.post('/api/roomba_cmd/_?cmd=' + mode, "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });

    var resetRoombaBtn = document.querySelector('#reset_roomba_btn');
    resetRoombaBtn.addEventListener('click', function() {
        $.post('/api/roomba_cmd/_?cmd=reset', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });

    var seekDockBtn = document.querySelector('#seek_dock_btn');
    seekDockBtn.addEventListener('click', function() {
        $.post('/api/roomba_cmd/_?cmd=seek_dock', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });


    var naviButton = document.querySelector('#navi-btn');
    naviButton.addEventListener('click', function() {
        $.post('/api/navi/', "", function(data, status) {
            if (data.Err != '') {
                console.log(data.Err);
                return
            }
        });
    });


});




// updateSpark updates the value of the tag and also the associated sparkline spark. 
// list is the buffer used to store the values.
function updateSpark(tag, list, spark, data, cnt) {
    $(tag).empty();
    $(tag).append(data);

    if (spark == "") {
        return
    }

    if (typeof list[tag] == "undefined") {
        list[tag] = [];
    }
    if (list[tag].length > 20) {
        list[tag].shift();
    }

    if (cnt % 3 == 0) {
        list[tag].push(data);
        $(spark).sparkline(list[tag]);
    }
}

// Handle websocket registrations and update Rover data panels.
$(document).ready(function() {

    var rb_batt_temp_list = [],
        rb_batt_volt_list = [],
        rb_batt_current_list = [],
        msgCount = 0,
        dataBuf = {},
        battCharge = 0,
        battPer = 0,
        total_velocity = 0,
        right_velocity = 0,
        left_velocity = 0,
        rb_main_curr = 0,
        rb_left_curr = 0,
        rb_right_curr = 0,
        rb_side_curr = 0;

    motor_velocity_chart = Morris.Bar({
        element: 'velocity_chart',
        data: [{
                y: 'Total',
                a: 100,
            },
            {
                y: 'Left',
                a: 75,
            },
            {
                y: 'Right',
                a: 50,
            },
        ],
        hideHover: 'auto',
        xkey: 'y',
        ykeys: ['a'],
        labels: ['mm/S', ]
    });


    motor_curr_chart = Morris.Bar({
        element: 'motor_current_chart',
        data: [{
                y: 'Left',
                curr: 100,
            },
            {
                y: 'Right',
                curr: 75,
            },
            {
                y: 'Main',
                curr: -50,
            },
            {
                y: 'Side',
                curr: 75,
            },
        ],
        hideHover: 'auto',
        xkey: 'y',
        ykeys: ['curr'],
        labels: ['mA', ]
    });


    ws = new WebSocket("ws://" + window.location.host + "/datastream");

    ws.onopen = function(evt) {
        $("#conn_spinner").show();
    }

    ws.onclose = function(evt) {
        $("#conn_spinner").hide();
        ws = null;
    }

    ws.onmessage = function(evt) {
        st = JSON.parse(evt.data);
        if (st.Err != "") {
            console.log(st.Err);
        }

        // Handle Roomba data feed.
        rbData = st.Roomba;
        for (var pktID in rbData) {
            pkt = rbData[pktID];
            switch (pktID) {
                case "7":
                    for (var i = 0; i < 4; i++) {
                        document.getElementById("rb_wheel_sensor").getElementsByTagName("tr")[i + 1].style.backgroundColor = "transparent";
                        if (pkt & (1 << i)) {
                            document.getElementById("rb_wheel_sensor").getElementsByTagName("tr")[i + 1].style.backgroundColor = "red";
                        }
                    }
                    break;

                case "8":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[5].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[5].style.backgroundColor = "transparent";
                    break;

                case "9":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[1].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[1].style.backgroundColor = "transparent";
                    break;

                case "10":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[2].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[2].style.backgroundColor = "transparent";
                    break;

                case "11":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[3].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[3].style.backgroundColor = "transparent";
                    break;

                case "12":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[4].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[4].style.backgroundColor = "transparent";
                    break;

                case "13":
                    // Virtual Wall
                    break;

                case "14":
                    for (var i = 0; i < 4; i++) {
                        document.getElementById("rb_overcurrent_sensor").getElementsByTagName("tr")[i + 1].style.backgroundColor = "transparent";
                    }
                    if (pkt & 1) {
                        document.getElementById("rb_overcurrent_sensor").getElementsByTagName("tr")[3].style.backgroundColor = "red";
                    }
                    if (pkt & 4) {
                        document.getElementById("rb_overcurrent_sensor").getElementsByTagName("tr")[4].style.backgroundColor = "red";
                    }
                    if (pkt & 8) {
                        document.getElementById("rb_overcurrent_sensor").getElementsByTagName("tr")[1].style.backgroundColor = "red";
                    }
                    if (pkt & 16) {
                        document.getElementById("rb_overcurrent_sensor").getElementsByTagName("tr")[2].style.backgroundColor = "red";
                    }
                    break;

                case "15":
                    // Dirt sensor.
                    break;

                case "16":
                    //unused.
                    break;

                case "17":
                    updateSpark("#rb_ir_omni", dataBuf, "", IR_CODE_NAMES[pkt], msgCount);
                    break;

                case "18":
                    // Buttons.
                    break;

                case "19":
                    //Distance.
                    break;

                case "20":
                    //Angle
                    break;

                case "21":
                    updateSpark("#rb_batt_charge_state", dataBuf, "", CHARGING_STATE[pkt], msgCount);
                    break;

                case "22":
                    updateSpark("#rb_batt_volt", dataBuf, ".rb_batt_volt_spark", pkt, msgCount);
                    break;

                case "23":
                    updateSpark("#rb_batt_current", dataBuf, ".rb_batt_current_spark", pkt, msgCount);
                    break;

                case "24":
                    updateSpark("#rb_batt_temp", dataBuf, ".rb_batt_temp_spark", pkt, msgCount);
                    break;

                case "25":
                    battCharge = pkt;
                    break;

                case "26":
                    battPer = Math.round(battCharge * 100 / pkt);
                    document.querySelector('.mdl-js-progress').MaterialProgress.setProgress(battPer);
                    $("#rb_batt_charge_tip").empty();
                    $("#rb_batt_charge_tip").append(battPer + "% " + battCharge + "/" + pkt);
                    updateSpark("#rb_batt_charge", [], "", battPer, msgCount);
                    break;

                case "27":
                    updateSpark("#rb_wall", [], "", pkt, msgCount);
                    break;

                case "28":
                    updateSpark("#rb_cliff_left", [], "", pkt, msgCount);
                    break;

                case "29":
                    updateSpark("#rb_cliff_front_left", [], "", pkt, msgCount);
                    break;

                case "30":
                    updateSpark("#rb_cliff_front_right", [], "", pkt, msgCount);
                    break;

                case "31":
                    updateSpark("#rb_cliff_right", [], "", pkt, msgCount);
                    break;

                case "32":
                    //unused.
                    break;

                case "33":
                    //unused.
                    break;

                case "34":
                    // Charge Source
                    break;

                case "35":
                    // OI mode.
                    updateSpark("#rb_oi_mode", [], "", OI_MODE[pkt], msgCount);
                    // Execute after an interval so it does not interfere with UI interaction.
                    setTimeout(function(pkt) {
                        switch (pkt) {
                            case 1:
                                $('#mode_passive_label')[0].MaterialRadio.check();
                                break;
                            case 2:
                                $('#mode_safe_label')[0].MaterialRadio.check();
                                break;
                            case 3:
                                $('#mode_full_label')[0].MaterialRadio.check();
                                break;
                            default:
                                console.log("Unknown mode");
                        }
                    }, 100, pkt);
                    break;

                case "36":
                    // Song Number.
                    break;

                case "37":
                    //Song Playing
                    break;

                case "38":
                    // OI stream num pkt.
                    break;

                case "39":
                    total_velocity = pkt;
                    break;

                case "40":
                    break;

                case "41":
                    right_velocity = pkt;
                    break;

                case "42":
                    left_velocity = pkt;
                    break;

                case "43":
                    // encoder counts.
                    break;

                case "44":
                    // encoder counts left.
                    break;

                case "45":
                    for (var i = 0; i < 6; i++) {
                        document.getElementById("rb_bump_sensor").getElementsByTagName("tr")[i + 1].style.backgroundColor = "transparent";
                    }
                    if (pkt & 1) {
                        document.getElementById("rb_bump_sensor").getElementsByTagName("tr")[1].style.backgroundColor = "red";
                    }
                    if (pkt & 2) {
                        document.getElementById("rb_bump_sensor").getElementsByTagName("tr")[2].style.backgroundColor = "red";
                    }
                    if (pkt & 4) {
                        document.getElementById("rb_bump_sensor").getElementsByTagName("tr")[3].style.backgroundColor = "red";
                    }
                    if (pkt & 8) {
                        document.getElementById("rb_bump_sensor").getElementsByTagName("tr")[4].style.backgroundColor = "red";
                    }
                    if (pkt & 16) {
                        document.getElementById("rb_bump_sensor").getElementsByTagName("tr")[5].style.backgroundColor = "red";
                    }
                    if (pkt & 32) {
                        document.getElementById("rb_bump_sensor").getElementsByTagName("tr")[6].style.backgroundColor = "red";
                    }
                    break;

                case "46":
                    updateSpark("#rb_bump_left", [], "", pkt, msgCount);
                    break;

                case "47":
                    updateSpark("#rb_bump_front_left", [], "", pkt, msgCount);
                    break;

                case "48":
                    updateSpark("#rb_bump_center_left", [], "", pkt, msgCount);
                    break;

                case "49":
                    updateSpark("#rb_bump_center_right", [], "", pkt, msgCount);
                    break;

                case "50":
                    updateSpark("#rb_bump_front_right", [], "", pkt, msgCount);
                    break;

                case "51":
                    updateSpark("#rb_bump_right", [], "", pkt, msgCount);
                    break;

                case "52":
                    updateSpark("#rb_ir_left", dataBuf, "", IR_CODE_NAMES[pkt], msgCount);
                    break;

                case "53":
                    updateSpark("#rb_ir_right", dataBuf, "", IR_CODE_NAMES[pkt], msgCount);
                    break;

                case "54":
                    rb_left_curr = pkt;
                    break;

                case "55":
                    rb_right_curr = pkt;
                    break;

                case "56":
                    rb_main_curr = pkt;

                case "57":
                    rb_side_curr = pkt;
                    break;
            }
        }
        msgCount++;

        motor_velocity_chart.setData([{
                y: 'Total',
                a: total_velocity,
            },
            {
                y: 'Left',
                a: left_velocity,
            },
            {
                y: 'Right',
                a: right_velocity,
            },
        ]);

        motor_curr_chart.setData([{
                y: 'Left',
                curr: rb_left_curr,
            },
            {
                y: 'Right',
                curr: rb_right_curr,
            },
            {
                y: 'Main',
                curr: rb_main_curr,
            },
            {
                y: 'Side',
                curr: rb_side_curr,
            },
        ]);

        ctrlData = st.Controller;
        for (var pktID in ctrlData) {
            pkt = ctrlData[pktID];
            switch (pktID) {
                case "0": // Temp.
                    updateSpark("#temp", dataBuf, ".temp_spark", parseFloat(pkt).toFixed(2), msgCount);
                    break;

                case "1": // Humidity.
                    updateSpark("#humidity", dataBuf, ".humidity_spark", pkt, msgCount);
                    break;

                case "2": // LDR.
                    updateSpark("#light", dataBuf, ".light_spark", pkt, msgCount);
                    break;

                case "3": // PIR.
                    updateSpark("#pir", dataBuf, ".pir_spark", pkt, msgCount);
                    break;

                case "4": // Heading.
                    updateSpark("#heading", [], "", parseFloat(pkt).toFixed(2), msgCount);
                    break;

                case "5": // Controller Volts.
                    updateSpark("#ctrl_volt", dataBuf, ".ctrl_volt_spark", parseFloat(pkt).toFixed(2), msgCount);
                    break;

            }
        }

        piData = st.Pi;
        for (var pktID in piData) {
            pkt = piData[pktID];
            switch (pktID) {
                case "0": //I2CBus State.
                    setTimeout(function(pkt) {
                        if (pkt == 1) {
                            $('#i2c_en_label')[0].MaterialSwitch.on();
                            return;
                        }
                        $('#i2c_en_label')[0].MaterialSwitch.off();
                        return;
                    }, 100, pkt);

                case "1": // AuxPower State.
                    setTimeout(function(pkt) {
                        if (pkt == 1) {
                            $('#aux_power_label')[0].MaterialSwitch.on();
                            return;
                        }
                        $('#aux_power_label')[0].MaterialSwitch.off();
                        return;
                    }, 100, pkt);

            }
        }

    }

    ws.onerror = function(evt) {
        print("ERROR: " + evt.data);
    }

});
