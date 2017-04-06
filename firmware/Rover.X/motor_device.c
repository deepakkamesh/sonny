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
void SpeedEncoderISR_M1(void);
void SpeedEncoderISR_M2(void);

void MotorTask(void) {
  uint8_t packet[PKT_SZ];
  static uint16_t rotation = 0;
  /*
    uint8_t currSt = 0;
    static uint8_t prevSt = 0;
    static uint32_t start, end = 0;

    // Pool rotations to speed encoder disk.
    if (active) {
      currSt = SE_M1_GetValue();
      if (currSt == 1 && prevSt == 0) { // Count on rising edge.
        start = TickGet();
      }
      if (currSt == 0 && prevSt == 1) {
        end = TickGet();
        if ((end - start) / TICK_MILLISECOND > 4) {
          se_m1_count++;
        }
      }
      prevSt = currSt;
    } */

  // Check if rotations are done.
  if ((se_m1_count >= rotation || se_m2_count >= rotation) && active) {
    MOTOR1_BWD_SetLow();
    MOTOR1_FWD_SetLow();
    MOTOR2_FWD_SetLow();
    MOTOR2_BWD_SetLow();
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

  switch (command) {
    case CMD_ON:
      se_m1_count = 0;
      se_m2_count = 0;
      MOTOR1_FWD_SetHigh();
      MOTOR1_BWD_SetLow();
      MOTOR2_FWD_SetHigh();
      MOTOR2_BWD_SetLow();
      SendDone(DEV_MOTOR);
      break;

    case CMD_FWD:
      if (CmdQ[DEV_MOTOR].size != 3) {
        // Send insufficient param error 
        SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
        break;
      }
      // Initialize counters.
      se_m1_count = 0;
      se_m2_count = 0;
      rotation = CmdQ[DEV_MOTOR].packet[1];
      rotation = rotation << 8 | CmdQ[DEV_MOTOR].packet[2];
      MOTOR1_FWD_SetHigh();
      MOTOR1_BWD_SetLow();
      MOTOR2_FWD_SetHigh();
      MOTOR2_BWD_SetLow();
      SendAck(DEV_MOTOR);
      active = true;
      break;

    case CMD_BWD:
      if (CmdQ[DEV_MOTOR].size != 3) {
        // Send insufficient param error 
        SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
        break;
      }
      se_m1_count = 0;
      se_m2_count = 0;
      rotation = CmdQ[DEV_MOTOR].packet[1];
      rotation = rotation << 8 | CmdQ[DEV_MOTOR].packet[2];
      MOTOR1_BWD_SetHigh();
      MOTOR1_FWD_SetLow();
      MOTOR2_BWD_SetHigh();
      MOTOR2_FWD_SetLow();
      SendAck(DEV_MOTOR);
      active = true;
      break;

    case CMD_OFF:
      MOTOR1_BWD_SetLow();
      MOTOR1_FWD_SetLow();
      MOTOR2_FWD_SetLow();
      MOTOR2_BWD_SetLow();
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
      if (dur > 4 && dur < 10) {
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
      if (dur > 4 && dur < 10) {
        se_m2_count++;
      }
    }
    prev = curr;
  }
}