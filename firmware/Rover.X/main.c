/**
  Generated Main Source File

  Company:
    Microchip Technology Inc.

  File Name:
    main.c

  Summary:
    This is the main file generated using MPLAB(c) Code Configurator

  Description:
    This header file provides implementations for driver APIs for all modules selected in the GUI.
    Generation Information :
        Product Revision  :  MPLAB(c) Code Configurator - 3.16
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
#include <stdbool.h>
#include "mcc_generated_files/mcc.h"
#include "protocol.h"
#include "host_controller.h"
#include "admin_device.h"
#include "led_device.h"
#include "servo_device.h"
#include "accel_device.h"
#include "ldr_device.h"
#include "batt_device.h"
#include "dht11_device.h"
#include "us020_device.h"
#include "lidar_device.h"
#include "tick.h"

/*
                         Main application
 */
void main(void) {
  // Initialize the device and subsystems.
  SYSTEM_Initialize();
  InitTicker();
  ServoInit();
  DHT11Init();
  HostControllerInit();
  
  // If using interrupts in PIC18 High/Low Priority Mode you need to enable the Global High and Low Interrupts
  // If using interrupts in PIC Mid-Range Compatibility Mode you need to enable the Global and Peripheral Interrupts
  // Use the following macros to:

  // Enable high priority global interrupts
  //INTERRUPT_GlobalInterruptHighEnable();

  // Enable low priority global interrupts.
  //INTERRUPT_GlobalInterruptLowEnable();

  // Disable high priority global interrupts
  //INTERRUPT_GlobalInterruptHighDisable();

  // Disable low priority global interrupts.
  //INTERRUPT_GlobalInterruptLowDisable();

  // Enable the Global Interrupts
  INTERRUPT_GlobalInterruptEnable();

  // Enable the Peripheral Interrupts
  INTERRUPT_PeripheralInterruptEnable();

  // Disable the Global Interrupts
  //INTERRUPT_GlobalInterruptDisable();

  // Disable the Peripheral Interrupts
  //INTERRUPT_PeripheralInterruptDisable();

  // Post Interrupt Enable Initialization.
   // LidarInit();

  __delay_ms(5);

  while (1) {
    AdminTask();
    LedTask();
    ServoTask();
    AccelTask();
    LDRTask();
    BattTask();
    DHT11Task();
    US020Task();
  //  LidarTask();
  }
}
/**
 End of File
 */