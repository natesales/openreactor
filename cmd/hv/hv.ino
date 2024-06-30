#include "DFRobot_GP8403.h"

#define LED_PIN LED_BUILTIN
#define TRIG_PIN 3
#define ADC_PIN 26
#define CURRENT_PIN 25

DFRobot_GP8403 dac(&Wire, 0x5F);

void setup() {
  pinMode(LED_PIN, OUTPUT);
  pinMode(TRIG_PIN, OUTPUT);
  pinMode(ADC_PIN, INPUT);
  pinMode(CURRENT_PIN, INPUT);
  analogReadResolution(12);
  digitalWrite(LED_PIN, LOW);
  digitalWrite(TRIG_PIN, LOW);

  Serial.begin(115200);

  uint8_t status = dac.begin();
  if (status != 0) {
    Serial.printf("I2C init error: %d", status);
    while (1) yield();
  }

  delay(100);
  dac.setDACOutRange(dac.eOutputRange5V);
  dac.setDACOutVoltage(0, 0);
  dac.setDACOutVoltage(0, 1);
  digitalWrite(LED_PIN, LOW);
}

void loop() {
  if (Serial.available() > 0) {
    char cmd = Serial.read();
    switch (cmd) {
      case 's':
        {
          int voltage = Serial.parseInt();
          if (voltage > 0) {
            digitalWrite(LED_PIN, HIGH);
          } else {
            digitalWrite(LED_PIN, LOW);
          }
          digitalWrite(TRIG_PIN, HIGH);
          dac.setDACOutVoltage(voltage, 0);
          Serial.printf("oks%d\n", voltage);
          digitalWrite(TRIG_PIN, LOW);
        }
        return;
      case 'r':
        {
          Serial.println(float(analogRead(ADC_PIN)) / 100);
        }
        return;
      case 'c':
        {
          float currentMa = float(analogRead(CURRENT_PIN)) / 1000;
          if (currentMa > 8) {
            dac.setDACOutVoltage(0, 0);
          }
          Serial.println(currentMa);
        }
        return;
      default:
        Serial.println("Unknown command");
        return;
    }
  }
}
