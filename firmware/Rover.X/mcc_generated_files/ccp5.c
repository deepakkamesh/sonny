/**
  CCP5 Generated Driver File

  @Company
    Microchip Technology Inc.

  @File Name
    ccp5.c

  @Summary
    This is the generated driver implementation file for the CCP5 driver using MPLAB(c) Code Configurator

  @Description
    This header file provides implementations for driver APIs for CCP5.
    Generation Information :
        Product Revision  :  MPLAB(c) Code Configurator - 3.15.0
        Device            :  PIC18F26K22
        Driver Version    :  2.00
    The generated drivers are tested against the following:
        Compiler          :  XC8 1.35
        MPLAB             :  MPLAB X 3.20
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
#include "ccp5.h"
#include "tmr5.h"

/**
  Section: COMPARE Module APIs
*/
static volatile  uint16_t pwm_on, pwm_off;

void CCP5_Initialize(void)
{
    // Set the CCP5 to the options selected in the User Interface
    
    // CCP5M Setoutput; DC5B 0; 
    CCP5CON = 0x08;
    
    // CCPR5L 0; 
    CCPR5L = 0x00;
    
    // CCPR5H 0; 
    CCPR5H = 0x00;
    
    // Selecting Timer 5
    CCPTMRS1bits.C5TSEL = 0x2;

    // Clear the CCP5 interrupt flag
    PIR4bits.CCP5IF = 0;
	
    // Enable the CCP5 interrupt
    PIE4bits.CCP5IE = 1;
    
    // Sane defaults for pwm.
    pwm_on = 1000;
    pwm_off = 9000;
}

void CCP5_SetCompareCount(uint16_t compareCount)
{
    CCP_PERIOD_REG_T module;
    
    // Write the 16-bit compare value
    module.ccpr5_16Bit = compareCount;
    
    CCPR5L = module.ccpr5l;
    CCPR5H = module.ccpr5h;
}

void CCP5_CompareISR(void)
{
    // Clear the CCP5 interrupt flag
    PIR4bits.CCP5IF = 0;
    
        // Reload timer with 0.
    TMR5_WriteTimer(0);
    if(CCP5CON == 8) {
        CCP5CON = 9;
        CCP5_SetCompareCount(pwm_on);
    }else {
        CCP5CON = 8;
        CCP5_SetCompareCount(pwm_off);
    }
}

void CCP5_SetOnOff(uint16_t on, uint16_t off) {
    pwm_on = on;
    pwm_off = off;
}
/**
 End of File
*/