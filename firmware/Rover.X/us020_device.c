#include "serial_controller.h"
#include "protocol.h"
#include "tick.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void US020Task(void) {


  if (CmdQ[DEV_US020].free) {
    // nothing to do
    return;
  }
  uint8_t command, packet[PKT_SZ];
  uint32_t ticks = 0;
  command = GetCommand(CmdQ[DEV_US020].packet[0]);

  switch (command) {
    case CMD_STATE:

      // Send trigger pulse.
      US_TRIG_SetHigh();
      __delay_us(10);
      US_TRIG_SetLow();

      // Read return echo pulse.
      // Wait for line to go high.
      ticks = TickGet();
      while (!US_ECHO_GetValue() && (TickGet() - ticks) / TICK_MILLISECOND < 25); // Max time to travel 400 cm

      // Read the time the line went high.
      ticks = TickGet();
      while (US_ECHO_GetValue() && (TickGet() - ticks) / TICK_MILLISECOND < 25);
      uint16_t dist = (TickGet() - ticks)*1000 / TICK_MILLISECOND / 58.3; 

      packet[0] = 0xC0 | DEV_US020; // Ack & Done.
      packet[1] = dist >> 8; 
      packet[2] = dist & 0xFF;
      SendPacket(packet, 3);
      break;

    default:
      SendError(DEV_US020, ERR_UNIMPLEMENTED);
      break;
  }

  CmdQ[DEV_US020].free = true;
}



