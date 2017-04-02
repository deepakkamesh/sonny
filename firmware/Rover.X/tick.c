#include "tick.h"
#include "mcc_generated_files/mcc.h"

// Global Variables.
static volatile uint32_t internalTicks = 0;
// 6-byte value to store Ticks.  Allows for use over longer periods of time.
static volatile uint8_t vTickReading[6];

// InitTicker initializes the ISR for Timer0.
void InitTicker(void) {
  TMR0_SetInterruptHandler(TMR0ISR);
}

// TMR0ISR is the timer0 ISR.
void TMR0ISR(void) {
  internalTicks++;
}

static void GetTickCopy(void) {
  do {
    INTCONbits.TMR0IE = 1; // Enable interrupt
    Nop();
    INTCONbits.TMR0IE = 0; // Disable interrupt
    vTickReading[0] = TMR0L;
    vTickReading[1] = TMR0H;
    *((uint16_t*) & vTickReading[2]) = internalTicks;
  } while (INTCONbits.TMR0IF);
  INTCONbits.TMR0IE = 1; // Enable interrupt
}

uint32_t TickGet(void) {
  uint32_t dw;

  GetTickCopy();
  ((uint8_t*) & dw)[0] = vTickReading[0]; // Note: This copy must be done one
  ((uint8_t*) & dw)[1] = vTickReading[1]; // byte at a time to prevent misaligned
  ((uint8_t*) & dw)[2] = vTickReading[2]; // memory reads, which will reset the PIC.
  ((uint8_t*) & dw)[3] = vTickReading[3];
  return dw;
}

uint32_t TickGetDev256(void) {
  uint32_t dw;

  GetTickCopy();
  ((uint8_t*) & dw)[0] = vTickReading[1]; // Note: This copy must be done one
  ((uint8_t*) & dw)[1] = vTickReading[2]; // byte at a time to prevent misaligned
  ((uint8_t*) & dw)[2] = vTickReading[3]; // memory reads, which will reset the PIC.
  ((uint8_t*) & dw)[3] = vTickReading[4];
  return dw;
}

uint32_t TickGetDev64K(void) {
  uint32_t dw;

  GetTickCopy();
  ((uint8_t*) & dw)[0] = vTickReading[2]; // Note: This copy must be done one
  ((uint8_t*) & dw)[1] = vTickReading[3]; // byte at a time to prevent misaligned
  ((uint8_t*) & dw)[2] = vTickReading[4]; // memory reads, which will reset the PIC.
  ((uint8_t*) & dw)[3] = vTickReading[5];
  return dw;
}
