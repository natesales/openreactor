#define INTERRUPT_PIN 15

uint32_t counts = 0;
uint16_t period = 1000;

void setup() {
  attachInterrupt(digitalPinToInterrupt(INTERRUPT_PIN), count, FALLING);
  Serial.begin(115200);
}

void count() {
  counts++;
}

void loop() {
  delay(period);
  detachInterrupt(digitalPinToInterrupt(INTERRUPT_PIN));
  Serial.println(counts);
  counts = 0;
  attachInterrupt(digitalPinToInterrupt(INTERRUPT_PIN), count, FALLING);
}
