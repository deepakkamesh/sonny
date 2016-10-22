/**
  Generated Pin Manager File

  Company:
    Microchip Technology Inc.

  File Name:
    pin_manager.c

  Summary:
    This is the Pin Manager file generated using MPLAB(c) Code Configurator

  Description:
    This header file provides implementations for pin APIs for all pins selected in the GUI.
    Generation Information :
        Product Revision  :  MPLAB(c) Code Configurator - 3.16
        Device            :  PIC18F26K22
        Driver Version    :  1.02
    The generated drivers are tested against the following:
        Compiler          :  XC8 1.35
        MPLAB             :  MPLAB X 3.20

    Copyright (c) 2013 - 2015 released Microchip Technology Inc.  All rights reserved.

    Microchip licenses to you the right to use, modify, copy and distribute
    Software only when embedded on a Microchip microcontroller or digital signal
    controller that is integrated into your product or third party product
    (pursuant to the sublicense terms in the accompanying license agreement).

    You should refer to the license agreement accompanying this Software for
    additional information regarding your rights and obligations.

    SOFTWARE AND DOCUMENTATION ARE PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND,
    EITHER EXPRESS OR IMPLIED, INCLUDING WITHOUT LIMITATION, ANY WARRANTY OF
    MERCHANTABILITY, TITLE, NON-INFRINGEMENT AND FITNESS FOR A PARTICULAR PURPOSE.
    IN NO EVENT SHALL MICROCHIP OR ITS LICENSORS BE LIABLE OR OBLIGATED UNDER
    CONTRACT, NEGLIGENCE, STRICT LIABILITY, CONTRIBUTION, BREACH OF WARRANTY, OR
    OTHER LEGAL EQUITABLE THEORY ANY DIRECT OR INDIRECT DAMAGES OR EXPENSES
    INCLUDING BUT NOT LIMITED TO ANY INCIDENTAL, SPECIAL, INDIRECT, PUNITIVE OR
    CONSEQUENTIAL DAMAGES, LOST PROFITS OR LOST DATA, COST OF PROCUREMENT OF
    SUBSTITUTE GOODS, TECHNOLOGY, SERVICES, OR ANY CLAIMS BY THIRD PARTIES
    (INCLUDING BUT NOT LIMITED TO ANY DEFENSE THEREOF), OR OTHER SIMILAR COSTS.

*/
#include <stdint.h>
#include <xc.h>
#include "pin_manager.h"

extern uint16_t se_m1_count,se_m2_count ;
void PIN_MANAGER_Initialize(void)
{
    LATB = 0x0;
    LATA = 0x0;
    LATC = 0x0;
    ANSELA = 0x2F;
    ANSELB = 0xC;
    ANSELC = 0x3C;
    TRISB = 0xFC;
    TRISC = 0xBC;
    WPUB = 0xF0;
    TRISA = 0x2F;

    INTCON2bits.nRBPU = 0x0;

    // interrupt on change for group IOCB - any
    IOCBbits.IOCB4 = 1; // Pin : RB4
    IOCBbits.IOCB5 = 1; // Pin : RB5

    INTCONbits.RBIE = 1; // Enable RBI interrupt 


}

void PIN_MANAGER_IOC(void)
{    
    // interrupt on change for group IOCB
    if(IOCBbits.IOCB4 == 1)
    {
        IOCB4_ISR();            
    }
    if(IOCBbits.IOCB5 == 1)
    {
        IOCB5_ISR();            
    }
}

/**
   IOCB4 Interrupt Service Routine
*/
void IOCB4_ISR(void) {

    // Add custom IOCB4 code
    if(IOCB4_InterruptHandler)
    {
        IOCB4_InterruptHandler();
    }
    IOCBbits.IOCB4 = 0;
}

/**
  Allows selecting an interrupt handler for IOCB4 at application runtime
*/
void IOCB4_SetInterruptHandler(void* InterruptHandler){
    IOCB4_InterruptHandler = InterruptHandler;
}

/**
  Default interrupt handler for IOCB4
*/
void IOCB4_DefaultInterruptHandler(void){
    // add your IOCB4 interrupt custom code
    // or set custom function using IOCB4_SetInterruptHandler()
    se_m1_count++;    LED1_SetHigh();

}
/**
   IOCB5 Interrupt Service Routine
*/
void IOCB5_ISR(void) {
LED1_SetHigh();
SE_M2_GetValue();
NOP();
    /*
    // Add custom IOCB5 code
    if(IOCB5_InterruptHandler)
    {
        IOCB5_InterruptHandler();
    }
    IOCBbits.IOCB5 = 0; */
}

/**
  Allows selecting an interrupt handler for IOCB5 at application runtime
*/
void IOCB5_SetInterruptHandler(void* InterruptHandler){
    IOCB5_InterruptHandler = InterruptHandler;
}

/**
  Default interrupt handler for IOCB5
*/
void IOCB5_DefaultInterruptHandler(void){
    // add your IOCB5 interrupt custom code
    // or set custom function using IOCB5_SetInterruptHandler()
    se_m2_count++;
    LED1_SetHigh();
}

/**
 End of File
*/
