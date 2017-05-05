#include "serial_controller.h"
#include "protocol.h"
#include "tick.h"
#include <stdlib.h>
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];
#define TIMEOUT 25

void US020Task(void) {
  static uint32_t ticks = 0;
  static uint16_t dist = 0;
  uint8_t packet[PKT_SZ];

  static enum {
    SEND_TRIGGER = 0,
    WAIT_HIGH,
    MEASURE_ECHO,
    SEND_RESPONSE,
    RESET,
    READY,
  } state = SEND_TRIGGER;

  // State Machine.
  switch (state) {
    case SEND_TRIGGER:
      US_TRIG_SetHigh();
      __delay_us(10);
      US_TRIG_SetLow();
      ticks = TickGet();
      state = WAIT_HIGH;
      break;

    case WAIT_HIGH:
      if ((TickGet() - ticks) / TICK_MILLISECOND > TIMEOUT) {
        SendError(DEV_US020, ERR_TIMEOUT);
        state = RESET;
        break;
      }
      if (!US_ECHO_GetValue()) {
        break;
      }
      ticks = TickGet();
      state = MEASURE_ECHO;
      break;

    case MEASURE_ECHO:
      if ((TickGet() - ticks) / TICK_MILLISECOND > TIMEOUT) {
        SendError(DEV_US020, ERR_TIMEOUT);
        state = RESET;
        break;
      }
      if (US_ECHO_GetValue()) {
        break;
      }
      dist = (TickGet() - ticks)*1000 / TICK_MILLISECOND / 58.3;
      state = SEND_RESPONSE;
      break;

    case SEND_RESPONSE:
      packet[0] = 0xC0 | DEV_US020; // Ack & Done.
      packet[1] = dist >> 8;
      packet[2] = dist & 0xFF;
      SendPacket(packet, 3);
      state = RESET;
      break;

    case RESET:
      CmdQ[DEV_US020].free = true;
      state = READY;
      break;

    case READY:
      break;
  }

  if (CmdQ[DEV_US020].free) {
    // nothing to do
    return;
  }

  uint8_t command = GetCommand(CmdQ[DEV_US020].packet[0]);

  switch (command) {
    case CMD_STATE:
      if (state != READY) {
        break;
      }
      state = SEND_TRIGGER;
      break;

    default:
      SendError(DEV_US020, ERR_UNIMPLEMENTED);
      break;
  }
}



