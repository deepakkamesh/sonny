#include <stdbool.h>
#include <stdlib.h>
#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"
#include "tick.h"

extern Queue CmdQ[MAX_DEVICES];

static volatile uint16_t se_m1_count = 0, se_m2_count = 0;
bool active = false, m1Done = false, m2Done = false;

#define M1_FWD_SetHigh STR1A=1
#define M1_FWD_SetLow STR1A=0
#define M1_BWD_SetHigh STR1B = 1
#define M1_BWD_SetLow STR1B = 0

#define M2_FWD_SetHigh STR2A = 1
#define M2_FWD_SetLow STR2A = 0
#define M2_BWD_SetHigh STR2B = 1
#define M2_BWD_SetLow STR2B = 0

enum uint8_t {
  RIGHT_SYNC = 0, // Async - Move both motors to turn.
  LEFT_SYNC, // Sync - Turn moving one motor.
  RIGHT_ASYNC,
  LEFT_ASYNC,
} dir = RIGHT_SYNC;

void SpeedEncoderISR_M1(void);
void SpeedEncoderISR_M2(void);
uint16_t getEncoderM1(void);
uint16_t getEncoderM2(void);

void MotorTask(void) {
  uint8_t packet[PKT_SZ];
  static uint16_t rotation = 0;

  // Turn off motors when # of rotations are met.
  if (active) {

    if (getEncoderM1() >= rotation) {
      M1_BWD_SetLow;
      M1_FWD_SetLow;
      m1Done = true;
    }
    if (getEncoderM2() >= rotation) {
      M2_BWD_SetLow;
      M2_FWD_SetLow;
      m2Done = true;
    }
    if (m1Done && m2Done) {
      packet[0] = 0x40 | DEV_MOTOR;
      packet[1] = se_m1_count >> 8;
      packet[2] = se_m1_count & 0xFF;
      packet[3] = se_m2_count >> 8;
      packet[4] = se_m2_count & 0xFF;
      SendPacket(packet, 5);
      m1Done = false, m2Done = false;
      active = false;
    }
  }

  if (CmdQ[DEV_MOTOR].free) {
    // nothing to do
    return;
  }

  uint8_t command = GetCommand(CmdQ[DEV_MOTOR].packet[0]);
  float dutyRatio = 0;
  uint16_t dutyValue = 0;


  switch (command) {
    case CMD_ON:
      if (CmdQ[DEV_MOTOR].size != 2) {
        SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
        break;
      }
      dutyRatio = (float) CmdQ[DEV_MOTOR].packet[1] / 100;
      dutyValue = 4 * (PR4 + 1) * dutyRatio;
      EPWM1_LoadDutyValue(dutyValue);
      EPWM2_LoadDutyValue(dutyValue);
      M1_FWD_SetHigh;
      M1_BWD_SetLow;
      M2_FWD_SetHigh;
      M2_BWD_SetLow;
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
      dutyRatio = (float) CmdQ[DEV_MOTOR].packet[3] / 100;
      dutyValue = 4 * (PR4 + 1) * dutyRatio;
      EPWM1_LoadDutyValue(dutyValue);
      EPWM2_LoadDutyValue(dutyValue);
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
      dutyRatio = (float) CmdQ[DEV_MOTOR].packet[3] / 100;
      dutyValue = 4 * (PR4 + 1) * dutyRatio;
      EPWM1_LoadDutyValue(dutyValue);
      EPWM2_LoadDutyValue(dutyValue);
      M1_BWD_SetHigh;
      M1_FWD_SetLow;
      M2_BWD_SetHigh;
      M2_FWD_SetLow;
      SendAck(DEV_MOTOR);
      active = true;
      break;

    case CMD_ROTATE:
      if (CmdQ[DEV_MOTOR].size != 5) {
        SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
        break;
      }
      se_m1_count = 0;
      se_m2_count = 0;
      rotation = CmdQ[DEV_MOTOR].packet[1];
      rotation = rotation << 8 | CmdQ[DEV_MOTOR].packet[2];
      dutyRatio = (float) CmdQ[DEV_MOTOR].packet[3] / 100;
      dutyValue = 4 * (PR4 + 1) * dutyRatio;
      EPWM1_LoadDutyValue(dutyValue);
      EPWM2_LoadDutyValue(dutyValue);
      dir = CmdQ[DEV_MOTOR].packet[4];
      switch (dir) {
        case RIGHT_SYNC:
          M2_BWD_SetHigh;
          M2_FWD_SetLow;
          M1_FWD_SetHigh;
          M1_BWD_SetLow;
          break;
        case LEFT_SYNC:
          M2_FWD_SetHigh;
          M2_BWD_SetLow;
          M1_BWD_SetHigh;
          M1_FWD_SetLow;
          break;
        case RIGHT_ASYNC:
          M2_BWD_SetHigh;
          M2_FWD_SetLow;
          M1_FWD_SetLow;
          M1_BWD_SetLow;
          break;
        case LEFT_ASYNC:
          M1_BWD_SetHigh;
          M1_FWD_SetLow;
          M2_FWD_SetLow;
          M2_BWD_SetLow;
          break;
      }
      SendAck(DEV_MOTOR);
      active = true;
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
      if (dur > 2 && dur < 20) {
        se_m1_count++;
      }
    }
    prev = curr;
  }
}

uint16_t getEncoderM1(void) {
  uint16_t val = 0;
  di();
  val = se_m1_count;
  ei();
  return val;
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
      if (dur > 2 && dur < 20) {
        se_m2_count++;
      }
    }
    prev = curr;
  }
}

uint16_t getEncoderM2(void) {
  uint16_t val = 0;
  di();
  val = se_m2_count;
  ei();
  return val;
}