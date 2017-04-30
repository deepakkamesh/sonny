/**
  CCP4 Generated Driver File

  @Company
    Microchip Technology Inc.

  @File Name
    ccp4.c

  @Summary
    This is the generated driver implementation file for the CCP4 driver using MPLAB(c) Code Configurator

  @Description
    This header file provides implementations for driver APIs for CCP4.
    Generation Information :
        Product Revision  :  MPLAB(c) Code Configurator - 4.15.1
        Device            :  PIC18F26K22
        Driver Version    :  2.00
    The generated drivers are tested against the following:
        Compiler          :  XC8 1.35
        MPLAB             :  MPLAB X 3.40
 */

/*
    (c) 2016 Microchip Technology Inc. and its subsidiaries. You may use this
    software and any derivatives exclusively with Microchip products.

    THIS SOFTWARE IS SUPPLIED BY MICROCHIP "AS IS". NO WARRANTIES, WHETHER
    EXPRESS, IMPLIED OR STATUTORY, APPLY TO THIS SOFTWARE, INCLUDING ANY IMPLIED
    WARRANTIES OF NON-INFRINGEMENT, MERCHANTABILITY, AND FITNESS FOR A
    PARTICULAR PURPOSE, OR ITS INTERACTION WITH MICROCHIP PRODUCTS, COMBINATION
    WITH ANY OTHER PRODUCTS, OR USE IN ANY APPLICATION.

    IN NO EVENT WILL MICROCHIP BE LIABLE FOR ANY INDIRECT, SPECIAL, PUNITIVE,
    INCIDENTAL OR CONSEQUENTIAL LOSS, DAMAGE, COST OR EXPENSE OF ANY KIND
    WHATSOEVER RELATED TO THE SOFTWARE, HOWEVER CAUSED, EVEN IF MICROCHIP HAS
    BEEN ADVISED OF THE POSSIBILITY OR THE DAMAGES ARE FORESEEABLE. TO THE
    FULLEST EXTENT ALLOWED BY LAW, MICROCHIP'S TOTAL LIABILITY ON ALL CLAIMS IN
    ANY WAY RELATED TO THIS SOFTWARE WILL NOT EXCEED THE AMOUNT OF FEES, IF ANY,
    THAT YOU HAVE PAID DIRECTLY TO MICROCHIP FOR THIS SOFTWARE.

    MICROCHIP PROVIDES THIS SOFTWARE CONDITIONALLY UPON YOUR ACCEPTANCE OF THESE
    TERMS.
 */

/**
  Section: Included Files
 */

#include <xc.h>
#include "ccp4.h"

/**
  Section: COMPARE Module APIs
 */
void (*CCP4_InterruptHandler)(void);

void CCP4_Initialize(void)
{
  // Set the CCP4 to the options selected in the User Interface

    // CCP4M Clearoutput; DC4B 0; 
    CCP4CON = 0x09;

  // CCPR4L 0; 
  CCPR4L = 0x00;

  // CCPR4H 0; 
  CCPR4H = 0x00;

  // Selecting Timer 3
  CCPTMRS1bits.C4TSEL = 0x1;

  // Clear the CCP4 interrupt flag
  PIR4bits.CCP4IF = 0;

  // Enable the CCP4 interrupt
  PIE4bits.CCP4IE = 1;
}

void CCP4_SetCompareCount(uint16_t compareCount)
{
  CCP_PERIOD_REG_T module;

  // Write the 16-bit compare value
  module.ccpr4_16Bit = compareCount;

  CCPR4L = module.ccpr4l;
  CCPR4H = module.ccpr4h;
}

void CCP4_CompareISR(void)
{
  // Clear the CCP4 interrupt flag
  PIR4bits.CCP4IF = 0;
  if (CCP4_InterruptHandler) {
    CCP4_InterruptHandler();
  }
}

void CCP4_SetInterruptHandler(void* InterruptHandler) {
  CCP4_InterruptHandler = InterruptHandler;
}
/**
 End of File
 */