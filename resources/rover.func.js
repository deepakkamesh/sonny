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
    //setInterval("RegularTasks()", 1000);
    var count = 0;
    var infoContainer = document.querySelector('#info-popup');

    var motorBackButton = document.querySelector('#motor-back');
    motorBackButton.addEventListener('click', function() {
        $("#rb_batt_temp").empty()
        $("#rb_batt_temp").append("w")

        //$("#servo-vert").empty()
        //$("#servo-vert").append(count++)
        document.getElementById("batt_metrics").getElementsByTagName("tr")[3].style.backgroundColor = "red";
        // alert("dd");
        // document.getElementById("batt_metrics").getElementsByTagName("tr")[3].style.backgroundColor = "transparent";


        $.post('/api/move/_?dir=back', "", function(data, status) {
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
        });
    });

    var motorForwardButton = document.querySelector('#motor-forward');
    motorForwardButton.addEventListener('click', function() {
        $.post('/api/move/_?dir=forward', "", function(data, status) {
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
        });
    });

    var motorLeftButton = document.querySelector('#motor-left');
    motorLeftButton.addEventListener('click', function() {
        $.post('/api/move/_?dir=left', "", function(data, status) {
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
        });
    });

    var motorRightButton = document.querySelector('#motor-right');
    motorRightButton.addEventListener('click', function() {
        $.post('/api/move/_?dir=right', "", function(data, status) {
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
        });
    });

    var driveVelSel = document.querySelector('#drive_velocity_sel');
    driveVelSel.addEventListener('click', function() {
        $("#drive_velocity_sel_disp").empty()
        $("#drive_velocity_sel_disp").append($('#drive_velocity_sel').val());
    });


    // Servo Controls.
    var servoDownButton = document.querySelector('#servo-down');
    servoDownButton.addEventListener('click', function() {
        $.post('/api/servorotate/_?dir=down', "", function(data, status) {
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#servo-up-tip").empty();
            $("#servo-up-tip").append(ret.Data['vert']);
            $("#servo-down-tip").empty();
            $("#servo-down-tip").append(ret.Data['vert']);
        });
    });

    var servoUpButton = document.querySelector('#servo-up');
    servoUpButton.addEventListener('click', function() {
        $.post('/api/servorotate/_?dir=up', "", function(data, status) {
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#servo-down-tip").empty();
            $("#servo-down-tip").append(ret.Data['vert']);
            $("#servo-up-tip").empty();
            $("#servo-up-tip").append(ret.Data['vert']);
        });
    });

    var servoLeftButton = document.querySelector('#servo-left');
    servoLeftButton.addEventListener('click', function() {
        $.post('/api/servorotate/_?dir=left', "", function(data, status) {
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#servo-left-tip").empty();
            $("#servo-left-tip").append(ret.Data['horiz']);
            $("#servo-right-tip").empty();
            $("#servo-right-tip").append(ret.Data['horiz']);
        });
    });

    var servoRightButton = document.querySelector('#servo-right');
    servoRightButton.addEventListener('click', function() {
        $.post('/api/servorotate/_?dir=right', "", function(data, status) {
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#servo-left-tip").empty();
            $("#servo-left-tip").append(ret.Data['horiz']);
            $("#servo-right-tip").empty();
            $("#servo-right-tip").append(ret.Data['horiz']);
        });
    });

    var servoAngleDeltaSel = document.querySelector('#servo_angle_step');
    servoAngleDeltaSel.addEventListener('click', function() {
        $("#servo_angle_step_disp").empty()
        $("#servo_angle_step_disp").append($('#servo_angle_step').val());
    });

    /*
       var LEDButton = document.querySelector('#led');
       LEDButton.addEventListener('click', function() {
           var action = '';
           if (document.getElementById('led').checked) {
               action = 'on';
           } else {
               action = 'off';
           }

           $.post('/api/ledon/_?cmd=' + action, "", function(data, status) {
               ret = JSON.parse(data)
               if (ret.Err != '') {
                   infoContainer.MaterialSnackbar.showSnackbar({
                       message: ret.Err
                   });
                   return
               }
           });
       });*/



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
        battPer = 0;

    ws = new WebSocket("ws://" + window.location.host + "/datastream");

    ws.onopen = function(evt) {}

    ws.onclose = function(evt) {
        ws = null;
    }

    ws.onmessage = function(evt) {
        st = JSON.parse(evt.data);
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

                case "8":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[5].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[5].style.backgroundColor = "transparent";

                case "9":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[1].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[1].style.backgroundColor = "transparent";

                case "10":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[2].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[2].style.backgroundColor = "transparent";

                case "11":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[3].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[3].style.backgroundColor = "transparent";

                case "12":
                    if (pkt == 1) {
                        document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[4].style.backgroundColor = "red";
                        continue;
                    }
                    document.getElementById("rb_cliff_sensor").getElementsByTagName("tr")[4].style.backgroundColor = "transparent"

                case "13":
                    // Virtual Wall

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

                case "15":
                    // Dirt sensor.

                case "16":
                    //unused.

                case "17":
                    updateSpark("#rb_ir_omni", dataBuf, "", IR_CODE_NAMES[pkt], msgCount);
                    break;

                case "18":
                    // Buttons.

                case "19":
                    //Distance.

                case "20":
                    //Angle

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

                    if (battCharge > pkt) continue; // TODO: Remove later. 


                    battPer = Math.round(battCharge * 100 / pkt);
                    document.querySelector('.mdl-js-progress').MaterialProgress.setProgress(battPer);

                    $("#rb_batt_charge_tip").empty();
                    $("#rb_batt_charge_tip").append(battPer + "% " + battCharge + "/" + pkt);


                case "28":
                    updateSpark("#rb_cliff_left", [], "", pkt, msgCount);

                case "29":
                    updateSpark("#rb_cliff_front_left", [], "", pkt, msgCount);

                case "30":
                    updateSpark("#rb_cliff_front_right", [], "", pkt, msgCount);

                case "31":
                    updateSpark("#rb_cliff_right", [], "", pkt, msgCount);

                case "52":
                    updateSpark("#rb_ir_left", dataBuf, "", IR_CODE_NAMES[pkt], msgCount);
                    break;

                case "53":
                    updateSpark("#rb_ir_right", dataBuf, "", IR_CODE_NAMES[pkt], msgCount);
                    break;

            }
        }
        msgCount++;
    }

    ws.onerror = function(evt) {
        print("ERROR: " + evt.data);
    }

});


/* Chart and graphs */
$(function() {

    /* CAtch keyboard */
    $(document).keyup(function(e) {
        if (e.which === 38) {
            alert("d")
            //up was pressed
        }
    });



    /* spark lines */
    $('.inlinesparkline').sparkline();

    var myvalues = [10, 8, 5, 7, 4, 4, 1];
    $('.dynamicsparkline').sparkline(myvalues);

    $('.dynamicbar').sparkline(myvalues, {
        type: 'bar',
        barColor: 'green'
    });

    $('.inlinebar').sparkline('html', {
        type: 'bar',
        barColor: 'red'
    });

    Morris.Bar({
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

    Morris.Bar({
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


});

// readRoverData opens a websocket
function readRoverData() {}

function RegularTasks() {
    var infoContainer = document.querySelector('#info-popup');

    // Ping controller.
    if (document.getElementById('master_en').checked) {
        $.post("/api/ping/", "", function(data, status) {
            ret = JSON.parse(data);
            $("#ping").empty();
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#ping").append(ret.Data);
        });
    }
    // Get distance.
    if (document.getElementById('distance_en').checked) {
        $.post("/api/distance/", "", function(data, status) {
            $("#distance").empty();
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#distance").append(ret.Data);
        });
    }

    // Get distance.
    if (document.getElementById('batt_en').checked) {
        $.post("/api/batt/", "", function(data, status) {
            $("#batt").empty();
            ret = JSON.parse(data)
            volts = Math.round(ret.Data * 1000) / 1000
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#batt").append(volts);
        });
    }

    // Get heading.
    if (document.getElementById('heading_en').checked) {
        $.post("/api/head/", "", function(data, status) {
            $("#heading").empty();
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            head = Math.round(ret.Data * 1000) / 1000

            $("#heading").append(head);
        });
    }

    // Get accel.
    if (document.getElementById('accel_en').checked) {
        $.post("/api/accel/", "", function(data, status) {
            $("#accel").empty();
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            x = Math.round(ret.Data[0] * 1000) / 1000
            y = Math.round(ret.Data[1] * 1000) / 1000
            z = Math.round(ret.Data[2] * 1000) / 1000
            $("#accel").append("X:" + x + " Y:" + y + " Z:" + z);
        });
    }

    // Get Temp & Humidity
    if (document.getElementById('temp_en').checked) {
        $.post("/api/temp/", "", function(data, status) {
            $("#temp").empty();
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#temp").append(" " + ret.Data[0] + "C " + ret.Data[1] + "%");
        });
    }

    // Get LDR
    if (document.getElementById('ldr_en').checked) {
        $.post("/api/ldr/", "", function(data, status) {
            $("#ldr").empty();
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#ldr").append(" " + ret.Data);
        });
    }

    // Get PIR
    if (document.getElementById('pir_en').checked) {
        $.post("/api/pir/", "", function(data, status) {
            $("#pir").empty();
            ret = JSON.parse(data)
            if (ret.Err != '') {
                infoContainer.MaterialSnackbar.showSnackbar({
                    message: ret.Err
                });
                return
            }
            $("#pir").append(ret.Data);
        });
    }
}