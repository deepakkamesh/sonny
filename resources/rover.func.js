$(document).ready(function() {
    //setInterval("RegularTasks()", 1000);
    var count = 0;
    var infoContainer = document.querySelector('#info-popup');

    var motorBackButton = document.querySelector('#motor-back');
    motorBackButton.addEventListener('click', function() {
        $("#rb_batt_temp").empty()
        $("#rb_batt_temp").append("w")

        $("#servo").empty()
        $("#servo").append(count++)
        document.getElementById("batt_metrics").getElementsByTagName("tr")[3].style.backgroundColor = "red";


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
            $("#servo-angle").empty();
            $("#servo-angle").append("Horiz:" + ret.Data['horiz'] + " Vert:" + ret.Data['vert']);
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
            $("#servo-angle").empty();
            $("#servo-angle").append("Horiz:" + ret.Data['horiz'] + " Vert:" + ret.Data['vert']);
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
            $("#servo-angle").empty();
            $("#servo-angle").append("Horiz:" + ret.Data['horiz'] + " Vert:" + ret.Data['vert']);
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
            $("#servo-angle").empty();
            $("#servo-angle").append("Horiz:" + ret.Data['horiz'] + " Vert:" + ret.Data['vert']);
        });
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

/* Chart and graphs */
$(function() {
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
