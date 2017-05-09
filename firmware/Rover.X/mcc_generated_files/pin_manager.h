/**
  @Generated Pin Manager Header File

  @Company:
    Microchip Technology Inc.

  @File Name:
    pin_manager.h

  @Summary:
    This is the Pin Manager file generated using MPLAB(c) Code Configurator

  @Description:
    This header file provides implementations for pin APIs for all pins selected in the GUI.
    Generation Information :
        Product Revision  :  MPLAB(c) Code Configurator - 4.15.1
        Device            :  PIC18F26K22
        Version           :  1.01
    The generated drivers are tested against the following:
        Compiler          :  XC8 1.35
        MPLAB             :  MPLAB X 3.40

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


#ifndef PIN_MANAGER_H
#define PIN_MANAGER_H

#define INPUT   1
#define OUTPUT  0

#define HIGH    1
#define LOW     0

#define ANALOG      1
#define DIGITAL     0

#define PULL_UP_ENABLED      1
#define PULL_UP_DISABLED     0

// get/set AX aliases
#define AX_TRIS               TRISAbits.TRISA0
#define AX_LAT                LATAbits.LATA0
#define AX_PORT               PORTAbits.RA0
#define AX_ANS                ANSELAbits.ANSA0
#define AX_SetHigh()            do { LATAbits.LATA0 = 1; } while(0)
#define AX_SetLow()             do { LATAbits.LATA0 = 0; } while(0)
#define AX_Toggle()             do { LATAbits.LATA0 = ~LATAbits.LATA0; } while(0)
#define AX_GetValue()           PORTAbits.RA0
#define AX_SetDigitalInput()    do { TRISAbits.TRISA0 = 1; } while(0)
#define AX_SetDigitalOutput()   do { TRISAbits.TRISA0 = 0; } while(0)
#define AX_SetAnalogMode()  do { ANSELAbits.ANSA0 = 1; } while(0)
#define AX_SetDigitalMode() do { ANSELAbits.ANSA0 = 0; } while(0)

// get/set AY aliases
#define AY_TRIS               TRISAbits.TRISA1
#define AY_LAT                LATAbits.LATA1
#define AY_PORT               PORTAbits.RA1
#define AY_ANS                ANSELAbits.ANSA1
#define AY_SetHigh()            do { LATAbits.LATA1 = 1; } while(0)
#define AY_SetLow()             do { LATAbits.LATA1 = 0; } while(0)
#define AY_Toggle()             do { LATAbits.LATA1 = ~LATAbits.LATA1; } while(0)
#define AY_GetValue()           PORTAbits.RA1
#define AY_SetDigitalInput()    do { TRISAbits.TRISA1 = 1; } while(0)
#define AY_SetDigitalOutput()   do { TRISAbits.TRISA1 = 0; } while(0)
#define AY_SetAnalogMode()  do { ANSELAbits.ANSA1 = 1; } while(0)
#define AY_SetDigitalMode() do { ANSELAbits.ANSA1 = 0; } while(0)

// get/set AZ aliases
#define AZ_TRIS               TRISAbits.TRISA2
#define AZ_LAT                LATAbits.LATA2
#define AZ_PORT               PORTAbits.RA2
#define AZ_ANS                ANSELAbits.ANSA2
#define AZ_SetHigh()            do { LATAbits.LATA2 = 1; } while(0)
#define AZ_SetLow()             do { LATAbits.LATA2 = 0; } while(0)
#define AZ_Toggle()             do { LATAbits.LATA2 = ~LATAbits.LATA2; } while(0)
#define AZ_GetValue()           PORTAbits.RA2
#define AZ_SetDigitalInput()    do { TRISAbits.TRISA2 = 1; } while(0)
#define AZ_SetDigitalOutput()   do { TRISAbits.TRISA2 = 0; } while(0)
#define AZ_SetAnalogMode()  do { ANSELAbits.ANSA2 = 1; } while(0)
#define AZ_SetDigitalMode() do { ANSELAbits.ANSA2 = 0; } while(0)

// get/set DHT11 aliases
#define DHT11_TRIS               TRISAbits.TRISA3
#define DHT11_LAT                LATAbits.LATA3
#define DHT11_PORT               PORTAbits.RA3
#define DHT11_ANS                ANSELAbits.ANSA3
#define DHT11_SetHigh()            do { LATAbits.LATA3 = 1; } while(0)
#define DHT11_SetLow()             do { LATAbits.LATA3 = 0; } while(0)
#define DHT11_Toggle()             do { LATAbits.LATA3 = ~LATAbits.LATA3; } while(0)
#define DHT11_GetValue()           PORTAbits.RA3
#define DHT11_SetDigitalInput()    do { TRISAbits.TRISA3 = 1; } while(0)
#define DHT11_SetDigitalOutput()   do { TRISAbits.TRISA3 = 0; } while(0)
#define DHT11_SetAnalogMode()  do { ANSELAbits.ANSA3 = 1; } while(0)
#define DHT11_SetDigitalMode() do { ANSELAbits.ANSA3 = 0; } while(0)

// get/set RA4 procedures
#define RA4_SetHigh()    do { LATAbits.LATA4 = 1; } while(0)
#define RA4_SetLow()   do { LATAbits.LATA4 = 0; } while(0)
#define RA4_Toggle()   do { LATAbits.LATA4 = ~LATAbits.LATA4; } while(0)
#define RA4_GetValue()         PORTAbits.RA4
#define RA4_SetDigitalInput()   do { TRISAbits.TRISA4 = 1; } while(0)
#define RA4_SetDigitalOutput()  do { TRISAbits.TRISA4 = 0; } while(0)

// get/set RB0 procedures
#define RB0_SetHigh()    do { LATBbits.LATB0 = 1; } while(0)
#define RB0_SetLow()   do { LATBbits.LATB0 = 0; } while(0)
#define RB0_Toggle()   do { LATBbits.LATB0 = ~LATBbits.LATB0; } while(0)
#define RB0_GetValue()         PORTBbits.RB0
#define RB0_SetDigitalInput()   do { TRISBbits.TRISB0 = 1; } while(0)
#define RB0_SetDigitalOutput()  do { TRISBbits.TRISB0 = 0; } while(0)
#define RB0_SetPullup()     do { WPUBbits.WPUB0 = 1; } while(0)
#define RB0_ResetPullup()   do { WPUBbits.WPUB0 = 0; } while(0)
#define RB0_SetAnalogMode() do { ANSELBbits.ANSB0 = 1; } while(0)
#define RB0_SetDigitalMode()do { ANSELBbits.ANSB0 = 0; } while(0)

// get/set LED1 aliases
#define LED1_TRIS               TRISBbits.TRISB1
#define LED1_LAT                LATBbits.LATB1
#define LED1_PORT               PORTBbits.RB1
#define LED1_WPU                WPUBbits.WPUB1
#define LED1_ANS                ANSELBbits.ANSB1
#define LED1_SetHigh()            do { LATBbits.LATB1 = 1; } while(0)
#define LED1_SetLow()             do { LATBbits.LATB1 = 0; } while(0)
#define LED1_Toggle()             do { LATBbits.LATB1 = ~LATBbits.LATB1; } while(0)
#define LED1_GetValue()           PORTBbits.RB1
#define LED1_SetDigitalInput()    do { TRISBbits.TRISB1 = 1; } while(0)
#define LED1_SetDigitalOutput()   do { TRISBbits.TRISB1 = 0; } while(0)
#define LED1_SetPullup()      do { WPUBbits.WPUB1 = 1; } while(0)
#define LED1_ResetPullup()    do { WPUBbits.WPUB1 = 0; } while(0)
#define LED1_SetAnalogMode()  do { ANSELBbits.ANSB1 = 1; } while(0)
#define LED1_SetDigitalMode() do { ANSELBbits.ANSB1 = 0; } while(0)

// get/set RB2 procedures
#define RB2_SetHigh()    do { LATBbits.LATB2 = 1; } while(0)
#define RB2_SetLow()   do { LATBbits.LATB2 = 0; } while(0)
#define RB2_Toggle()   do { LATBbits.LATB2 = ~LATBbits.LATB2; } while(0)
#define RB2_GetValue()         PORTBbits.RB2
#define RB2_SetDigitalInput()   do { TRISBbits.TRISB2 = 1; } while(0)
#define RB2_SetDigitalOutput()  do { TRISBbits.TRISB2 = 0; } while(0)
#define RB2_SetPullup()     do { WPUBbits.WPUB2 = 1; } while(0)
#define RB2_ResetPullup()   do { WPUBbits.WPUB2 = 0; } while(0)
#define RB2_SetAnalogMode() do { ANSELBbits.ANSB2 = 1; } while(0)
#define RB2_SetDigitalMode()do { ANSELBbits.ANSB2 = 0; } while(0)

// get/set SE_M1 aliases
#define SE_M1_TRIS               TRISBbits.TRISB4
#define SE_M1_LAT                LATBbits.LATB4
#define SE_M1_PORT               PORTBbits.RB4
#define SE_M1_WPU                WPUBbits.WPUB4
#define SE_M1_ANS                ANSELBbits.ANSB4
#define SE_M1_SetHigh()            do { LATBbits.LATB4 = 1; } while(0)
#define SE_M1_SetLow()             do { LATBbits.LATB4 = 0; } while(0)
#define SE_M1_Toggle()             do { LATBbits.LATB4 = ~LATBbits.LATB4; } while(0)
#define SE_M1_GetValue()           PORTBbits.RB4
#define SE_M1_SetDigitalInput()    do { TRISBbits.TRISB4 = 1; } while(0)
#define SE_M1_SetDigitalOutput()   do { TRISBbits.TRISB4 = 0; } while(0)
#define SE_M1_SetPullup()      do { WPUBbits.WPUB4 = 1; } while(0)
#define SE_M1_ResetPullup()    do { WPUBbits.WPUB4 = 0; } while(0)
#define SE_M1_SetAnalogMode()  do { ANSELBbits.ANSB4 = 1; } while(0)
#define SE_M1_SetDigitalMode() do { ANSELBbits.ANSB4 = 0; } while(0)

// get/set SE_M2 aliases
#define SE_M2_TRIS               TRISBbits.TRISB5
#define SE_M2_LAT                LATBbits.LATB5
#define SE_M2_PORT               PORTBbits.RB5
#define SE_M2_WPU                WPUBbits.WPUB5
#define SE_M2_ANS                ANSELBbits.ANSB5
#define SE_M2_SetHigh()            do { LATBbits.LATB5 = 1; } while(0)
#define SE_M2_SetLow()             do { LATBbits.LATB5 = 0; } while(0)
#define SE_M2_Toggle()             do { LATBbits.LATB5 = ~LATBbits.LATB5; } while(0)
#define SE_M2_GetValue()           PORTBbits.RB5
#define SE_M2_SetDigitalInput()    do { TRISBbits.TRISB5 = 1; } while(0)
#define SE_M2_SetDigitalOutput()   do { TRISBbits.TRISB5 = 0; } while(0)
#define SE_M2_SetPullup()      do { WPUBbits.WPUB5 = 1; } while(0)
#define SE_M2_ResetPullup()    do { WPUBbits.WPUB5 = 0; } while(0)
#define SE_M2_SetAnalogMode()  do { ANSELBbits.ANSB5 = 1; } while(0)
#define SE_M2_SetDigitalMode() do { ANSELBbits.ANSB5 = 0; } while(0)

// get/set RC0 procedures
#define RC0_SetHigh()    do { LATCbits.LATC0 = 1; } while(0)
#define RC0_SetLow()   do { LATCbits.LATC0 = 0; } while(0)
#define RC0_Toggle()   do { LATCbits.LATC0 = ~LATCbits.LATC0; } while(0)
#define RC0_GetValue()         PORTCbits.RC0
#define RC0_SetDigitalInput()   do { TRISCbits.TRISC0 = 1; } while(0)
#define RC0_SetDigitalOutput()  do { TRISCbits.TRISC0 = 0; } while(0)

// get/set RC1 procedures
#define RC1_SetHigh()    do { LATCbits.LATC1 = 1; } while(0)
#define RC1_SetLow()   do { LATCbits.LATC1 = 0; } while(0)
#define RC1_Toggle()   do { LATCbits.LATC1 = ~LATCbits.LATC1; } while(0)
#define RC1_GetValue()         PORTCbits.RC1
#define RC1_SetDigitalInput()   do { TRISCbits.TRISC1 = 1; } while(0)
#define RC1_SetDigitalOutput()  do { TRISCbits.TRISC1 = 0; } while(0)

// get/set RC2 procedures
#define RC2_SetHigh()    do { LATCbits.LATC2 = 1; } while(0)
#define RC2_SetLow()   do { LATCbits.LATC2 = 0; } while(0)
#define RC2_Toggle()   do { LATCbits.LATC2 = ~LATCbits.LATC2; } while(0)
#define RC2_GetValue()         PORTCbits.RC2
#define RC2_SetDigitalInput()   do { TRISCbits.TRISC2 = 1; } while(0)
#define RC2_SetDigitalOutput()  do { TRISCbits.TRISC2 = 0; } while(0)
#define RC2_SetAnalogMode() do { ANSELCbits.ANSC2 = 1; } while(0)
#define RC2_SetDigitalMode()do { ANSELCbits.ANSC2 = 0; } while(0)

// get/set LDR aliases
#define LDR_TRIS               TRISCbits.TRISC3
#define LDR_LAT                LATCbits.LATC3
#define LDR_PORT               PORTCbits.RC3
#define LDR_ANS                ANSELCbits.ANSC3
#define LDR_SetHigh()            do { LATCbits.LATC3 = 1; } while(0)
#define LDR_SetLow()             do { LATCbits.LATC3 = 0; } while(0)
#define LDR_Toggle()             do { LATCbits.LATC3 = ~LATCbits.LATC3; } while(0)
#define LDR_GetValue()           PORTCbits.RC3
#define LDR_SetDigitalInput()    do { TRISCbits.TRISC3 = 1; } while(0)
#define LDR_SetDigitalOutput()   do { TRISCbits.TRISC3 = 0; } while(0)
#define LDR_SetAnalogMode()  do { ANSELCbits.ANSC3 = 1; } while(0)
#define LDR_SetDigitalMode() do { ANSELCbits.ANSC3 = 0; } while(0)

// get/set US_TRIG aliases
#define US_TRIG_TRIS               TRISCbits.TRISC4
#define US_TRIG_LAT                LATCbits.LATC4
#define US_TRIG_PORT               PORTCbits.RC4
#define US_TRIG_ANS                ANSELCbits.ANSC4
#define US_TRIG_SetHigh()            do { LATCbits.LATC4 = 1; } while(0)
#define US_TRIG_SetLow()             do { LATCbits.LATC4 = 0; } while(0)
#define US_TRIG_Toggle()             do { LATCbits.LATC4 = ~LATCbits.LATC4; } while(0)
#define US_TRIG_GetValue()           PORTCbits.RC4
#define US_TRIG_SetDigitalInput()    do { TRISCbits.TRISC4 = 1; } while(0)
#define US_TRIG_SetDigitalOutput()   do { TRISCbits.TRISC4 = 0; } while(0)
#define US_TRIG_SetAnalogMode()  do { ANSELCbits.ANSC4 = 1; } while(0)
#define US_TRIG_SetDigitalMode() do { ANSELCbits.ANSC4 = 0; } while(0)

// get/set US_ECHO aliases
#define US_ECHO_TRIS               TRISCbits.TRISC5
#define US_ECHO_LAT                LATCbits.LATC5
#define US_ECHO_PORT               PORTCbits.RC5
#define US_ECHO_ANS                ANSELCbits.ANSC5
#define US_ECHO_SetHigh()            do { LATCbits.LATC5 = 1; } while(0)
#define US_ECHO_SetLow()             do { LATCbits.LATC5 = 0; } while(0)
#define US_ECHO_Toggle()             do { LATCbits.LATC5 = ~LATCbits.LATC5; } while(0)
#define US_ECHO_GetValue()           PORTCbits.RC5
#define US_ECHO_SetDigitalInput()    do { TRISCbits.TRISC5 = 1; } while(0)
#define US_ECHO_SetDigitalOutput()   do { TRISCbits.TRISC5 = 0; } while(0)
#define US_ECHO_SetAnalogMode()  do { ANSELCbits.ANSC5 = 1; } while(0)
#define US_ECHO_SetDigitalMode() do { ANSELCbits.ANSC5 = 0; } while(0)

// get/set RC6 procedures
#define RC6_SetHigh()    do { LATCbits.LATC6 = 1; } while(0)
#define RC6_SetLow()   do { LATCbits.LATC6 = 0; } while(0)
#define RC6_Toggle()   do { LATCbits.LATC6 = ~LATCbits.LATC6; } while(0)
#define RC6_GetValue()         PORTCbits.RC6
#define RC6_SetDigitalInput()   do { TRISCbits.TRISC6 = 1; } while(0)
#define RC6_SetDigitalOutput()  do { TRISCbits.TRISC6 = 0; } while(0)
#define RC6_SetAnalogMode() do { ANSELCbits.ANSC6 = 1; } while(0)
#define RC6_SetDigitalMode()do { ANSELCbits.ANSC6 = 0; } while(0)

// get/set RC7 procedures
#define RC7_SetHigh()    do { LATCbits.LATC7 = 1; } while(0)
#define RC7_SetLow()   do { LATCbits.LATC7 = 0; } while(0)
#define RC7_Toggle()   do { LATCbits.LATC7 = ~LATCbits.LATC7; } while(0)
#define RC7_GetValue()         PORTCbits.RC7
#define RC7_SetDigitalInput()   do { TRISCbits.TRISC7 = 1; } while(0)
#define RC7_SetDigitalOutput()  do { TRISCbits.TRISC7 = 0; } while(0)
#define RC7_SetAnalogMode() do { ANSELCbits.ANSC7 = 1; } while(0)
#define RC7_SetDigitalMode()do { ANSELCbits.ANSC7 = 0; } while(0)

/**
   @Param
    none
   @Returns
    none
   @Description
    GPIO and peripheral I/O initialization
   @Example
    PIN_MANAGER_Initialize();
 */
void PIN_MANAGER_Initialize (void);

/**
 * @Param
    none
 * @Returns
    none
 * @Description
    Interrupt on Change Handling routine
 * @Example
    PIN_MANAGER_IOC();
 */
void PIN_MANAGER_IOC(void);


/**
 * @Param
    none
 * @Returns
    none
 * @Description
    Interrupt on Change Handler for the IOCB4 pin functionality
 * @Example
    IOCB4_ISR();
 */
void IOCB4_ISR(void);

/**
  @Summary
    Interrupt Handler Setter for IOCB4 pin interrupt-on-change functionality

  @Description
    Allows selecting an interrupt handler for IOCB4 at application runtime
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    InterruptHandler function pointer.

  @Example
    PIN_MANAGER_Initialize();
    IOCB4_SetInterruptHandler(MyInterruptHandler);

*/
void IOCB4_SetInterruptHandler(void* InterruptHandler);

/**
  @Summary
    Dynamic Interrupt Handler for IOCB4 pin

  @Description
    This is a dynamic interrupt handler to be used together with the IOCB4_SetInterruptHandler() method.
    This handler is called every time the IOCB4 ISR is executed and allows any function to be registered at runtime.
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    None.

  @Example
    PIN_MANAGER_Initialize();
    IOCB4_SetInterruptHandler(IOCB4_InterruptHandler);

*/
extern void (*IOCB4_InterruptHandler)(void);

/**
  @Summary
    Default Interrupt Handler for IOCB4 pin

  @Description
    This is a predefined interrupt handler to be used together with the IOCB4_SetInterruptHandler() method.
    This handler is called every time the IOCB4 ISR is executed. 
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    None.

  @Example
    PIN_MANAGER_Initialize();
    IOCB4_SetInterruptHandler(IOCB4_DefaultInterruptHandler);

*/
void IOCB4_DefaultInterruptHandler(void);


/**
 * @Param
    none
 * @Returns
    none
 * @Description
    Interrupt on Change Handler for the IOCB5 pin functionality
 * @Example
    IOCB5_ISR();
 */
void IOCB5_ISR(void);

/**
  @Summary
    Interrupt Handler Setter for IOCB5 pin interrupt-on-change functionality

  @Description
    Allows selecting an interrupt handler for IOCB5 at application runtime
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    InterruptHandler function pointer.

  @Example
    PIN_MANAGER_Initialize();
    IOCB5_SetInterruptHandler(MyInterruptHandler);

*/
void IOCB5_SetInterruptHandler(void* InterruptHandler);

/**
  @Summary
    Dynamic Interrupt Handler for IOCB5 pin

  @Description
    This is a dynamic interrupt handler to be used together with the IOCB5_SetInterruptHandler() method.
    This handler is called every time the IOCB5 ISR is executed and allows any function to be registered at runtime.
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    None.

  @Example
    PIN_MANAGER_Initialize();
    IOCB5_SetInterruptHandler(IOCB5_InterruptHandler);

*/
extern void (*IOCB5_InterruptHandler)(void);

/**
  @Summary
    Default Interrupt Handler for IOCB5 pin

  @Description
    This is a predefined interrupt handler to be used together with the IOCB5_SetInterruptHandler() method.
    This handler is called every time the IOCB5 ISR is executed. 
    
  @Preconditions
    Pin Manager intializer called

  @Returns
    None.

  @Param
    None.

  @Example
    PIN_MANAGER_Initialize();
    IOCB5_SetInterruptHandler(IOCB5_DefaultInterruptHandler);

*/
void IOCB5_DefaultInterruptHandler(void);



#endif // PIN_MANAGER_H
/**
 End of File
*/