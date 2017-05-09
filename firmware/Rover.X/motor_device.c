#include <stdbool.h>
#include <stdlib.h>
#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"
#include "tick.h"

extern Queue CmdQ[MAX_DEVICES];

static volatile uint16_t se_m1_count = 0, se_m2_count = 0;
bool active = false;

#define M1_FWD_SetHigh STR1A=1
#define M1_FWD_SetLow STR1A=0
#define M1_BWD_SetHigh STR1B = 1
#define M1_BWD_SetLow STR1B = 0

#define M2_FWD_SetHigh STR2A = 1
#define M2_FWD_SetLow STR2A = 0
#define M2_BWD_SetHigh STR2B = 1
#define M2_BWD_SetLow STR2B = 0

enum uint8_t {
  RIGHT_SYNC = 0,
  LEFT_SYNC,
  RIGHT_ASYNC,
  LEFT_ASYNC,
} dir = RIGHT_SYNC;

void SpeedEncoderISR_M1(void);
void SpeedEncoderISR_M2(void);

void MotorTask(void) {
  uint8_t packet[PKT_SZ];
  static uint16_t rotation = 0;

  // Check if rotations are done.
  if ((se_m1_count >= rotation || se_m2_count >= rotation) && active) {
    M1_BWD_SetLow;
    M1_FWD_SetLow;
    M2_FWD_SetLow;
    M2_BWD_SetLow;
    packet[0] = 0x40 | DEV_MOTOR;
    packet[1] = se_m1_count >> 8;
    packet[2] = se_m1_count & 0xFF;
    packet[3] = se_m2_count >> 8;
    packet[4] = se_m2_count & 0xFF;
    SendPacket(packet, 5);
    active = false;
  }

  if (CmdQ[DEV_MOTOR].free) {
    // nothing to do
    return;
  }

  uint8_t command = GetCommand(CmdQ[DEV_MOTOR].packet[0]);
  float dutyRatio = 0;
  uint16_t dutyValue1, dutyValue2 = 0;


  switch (command) {
    case CMD_ON:
      if (CmdQ[DEV_MOTOR].size != 2) {
        SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
        break;
      }
      dutyRatio = (float) CmdQ[DEV_MOTOR].packet[1] / 100;
      dutyValue1 = 4 * (PR4 + 1) * dutyRatio;
      dutyValue2 = 4 * (PR6 + 1) * dutyRatio;

      M1_FWD_SetHigh;
      M1_BWD_SetLow;
      M2_FWD_SetHigh;
      M2_BWD_SetLow;
      EPWM1_LoadDutyValue(dutyValue1);
      EPWM2_LoadDutyValue(dutyValue2);
      SendAckDone(DEV_MOTOR);
      active = false;
      break;

    case CMD_FWD:
      if (CmdQ[DEV_MOTOR].size != 4) {
        SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
        break;
      }
      se_m1_count = 0;
      se_m2_count = 0;
      rotation = CmdQ[DEV_MOTOR].packet[1];
      rotation = rotation << 8 | CmdQ[DEV_MOTOR].packet[2];
      dutyRatio = (float) CmdQ[DEV_MOTOR].packet[1] / 100;
      dutyValue1 = 4 * (PR4 + 1) * dutyRatio;
      dutyValue2 = 4 * (PR6 + 1) * dutyRatio;
      EPWM1_LoadDutyValue(dutyValue1);
      EPWM2_LoadDutyValue(dutyValue2);
      M1_FWD_SetHigh;
      M1_BWD_SetLow;
      M2_FWD_SetHigh;
      M2_BWD_SetLow;
      SendAck(DEV_MOTOR);
      active = true;
      break;

    case CMD_BWD:
      if (CmdQ[DEV_MOTOR].size != 4) {
        SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
        break;
      }
      se_m1_count = 0;
      se_m2_count = 0;
      rotation = CmdQ[DEV_MOTOR].packet[1];
      rotation = rotation << 8 | CmdQ[DEV_MOTOR].packet[2];
      dutyRatio = (float) CmdQ[DEV_MOTOR].packet[1] / 100;
      dutyValue1 = 4 * (PR4 + 1) * dutyRatio;
      dutyValue2 = 4 * (PR6 + 1) * dutyRatio;
      EPWM1_LoadDutyValue(dutyValue1);
      EPWM2_LoadDutyValue(dutyValue2);
      M1_BWD_SetHigh;
      M1_FWD_SetLow;
      M2_BWD_SetHigh;
      M2_FWD_SetLow;
      SendAck(DEV_MOTOR);
      active = true;
      break;

    case CMD_ROTATE:

      break;


    case CMD_OFF:
      M1_BWD_SetLow;
      M1_FWD_SetLow;
      M2_FWD_SetLow;
      M2_BWD_SetLow;
      SendAckDone(DEV_MOTOR);
      active = false;
      break;

    case CMD_STATE:
      packet[0] = 0xC0 | DEV_MOTOR;
      packet[1] = se_m1_count >> 8;
      packet[2] = se_m1_count & 0xFF;
      packet[3] = se_m2_count >> 8;
      packet[4] = se_m2_count & 0xFF;
      SendPacket(packet, 5);
      break;

    default:
      SendError(DEV_MOTOR, ERR_UNIMPLEMENTED);
      break;
  }

  CmdQ[DEV_MOTOR].free = true;
}

void MotorInit(void) {
  IOCB4_SetInterruptHandler(SpeedEncoderISR_M1);
  IOCB5_SetInterruptHandler(SpeedEncoderISR_M2);

  // Turn off motors otherwise the default pwm is 50% duty cycle.
  M1_BWD_SetLow;
  M1_FWD_SetLow;
  M2_FWD_SetLow;
  M2_BWD_SetLow;
}

void SpeedEncoderISR_M1(void) {
  static char prev = 0;
  static char curr = 0;
  static uint32_t end = 0;
  static uint32_t start = 0;

  curr = SE_M1_GetValue();
  NOP();

  // Pool rotations to speed encoder disk.
  if (active) {
    if (curr == 1 && prev == 0) {
      start = TickGet();
    }
    if (curr == 0 && prev == 1) {
      end = TickGet();
      uint32_t dur = (end - start) / TICK_MILLISECOND;
      if (dur > 2 && dur < 10) {
        se_m1_count++;
      }
    }
    prev = curr;
  }
}

void SpeedEncoderISR_M2(void) {
  static char prev = 0;
  static char curr = 0;
  static uint32_t end = 0;
  static uint32_t start = 0;

  curr = SE_M2_GetValue();
  NOP();

  // Pool rotations to speed encoder disk.
  if (active) {

    if (curr == 1 && prev == 0) {
      start = TickGet();
    }
    if (curr == 0 && prev == 1) {
      end = TickGet();
      uint32_t dur = (end - start) / TICK_MILLISECOND;
      if (dur > 2 && dur < 10) {
        se_m2_count++;
      }
    }
    prev = curr;
  }
}