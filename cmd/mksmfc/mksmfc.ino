#define LED_PIN LED_BUILTIN
#define TRIG_PIN 27
#define ADC_PIN 26
#define SETPOINT_PIN 25

void setup() {
  Serial.begin(115200);

  pinMode(LED_PIN, OUTPUT);
  pinMode(TRIG_PIN, OUTPUT);
  pinMode(SETPOINT_PIN, OUTPUT);
  pinMode(ADC_PIN, INPUT);
  analogReadResolution(12);
  analogWriteResolution(16);

  digitalWrite(LED_PIN, LOW);
  analogWrite(LED_PIN, 0);
}

void loop() {
  if (Serial.available() > 0) {
    char cmd = Serial.read();
    switch (cmd) {
      case 's':
        {
          int voltage = Serial.parseInt();
          if (voltage > 0) {
            digitalWrite(TRIG_PIN, HIGH);
          } else {
            digitalWrite(TRIG_PIN, LOW);
          }
          analogWrite(LED_PIN, voltage);
          Serial.printf("oks%d\n", voltage);
        }
        return;
      case 'r':
        {
          Serial.println(float(analogRead(ADC_PIN)) / 100);
        }
        return;
      default:
        Serial.println("Unknown command");
        return;
    }
  }
}
