#include <stdbool.h>
#include <stdlib.h>
#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "servo_device.h"
#include "mcc_generated_files/mcc.h"


extern Queue CmdQ[MAX_DEVICES];
volatile static uint16_t pwm4_on, pwm4_off, pwm5_on, pwm5_off;
void CCP4_ISR(void);
void CCP5_ISR(void);

void ServoTask(void) {
  if (CmdQ[DEV_SERVO].free) {
    // nothing to do
    return;
  }
  uint8_t command, servo;
  command = GetCommand(CmdQ[DEV_SERVO].packet[0]);
  uint16_t on, period;
  switch (command) {
    case CMD_ROTATE:
      on = 1000; //default duration. Center.
      // Load on time. TODO set limits.
      if (CmdQ[DEV_SERVO].size != 6) {
        // Send insufficient param error 
        SendError(DEV_SERVO, ERR_INSUFFICENT_PARAMS);
        break;
      }
      on = CmdQ[DEV_SERVO].packet[1];
      on = on << 8 | CmdQ[DEV_SERVO].packet[2];
      period = CmdQ[DEV_SERVO].packet[3];
      period = period << 8 | CmdQ[DEV_SERVO].packet[4];
      servo = CmdQ[DEV_SERVO].packet[5]; // Servo Select.
      switch (servo) {
        case 0x1:
          pwm4_on = on;
          pwm4_off = period - on;
          break;
        case 0x2:
          pwm5_on = on;
          pwm5_off = period - on;
          break;
      }
      SendAckDone(DEV_SERVO);
      break;

    default:
      SendError(DEV_SERVO, ERR_UNIMPLEMENTED);
      break;

  }
  CmdQ[DEV_SERVO].free = true;
}

void ServoInit(void) {

  pwm4_on = 3000;
  pwm4_off = 37000;
  pwm5_on = 3000;
  pwm5_off = 37000;
  CCP4_SetInterruptHandler(CCP4_ISR);
  CCP5_SetInterruptHandler(CCP5_ISR);
}

void CCP5_ISR(void) {
  // Reload timer with 0.
  TMR5_WriteTimer(0);
  if (CCP5CON == 8) {
    CCP5CON = 9;
    CCP5_SetCompareCount(pwm5_on);
  } else {
    CCP5CON = 8;
    CCP5_SetCompareCount(pwm5_off);
  }
}

void CCP4_ISR(void) {
  // Reload timer with 0.
  TMR3_WriteTimer(0);
  if (CCP4CON == 8) {
    CCP4CON = 9;
    CCP4_SetCompareCount(pwm4_on);
  } else {
    CCP4CON = 8;
    CCP4_SetCompareCount(pwm4_off);
  }
}
