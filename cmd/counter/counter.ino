#define INTERRUPT_PIN 15

int counts = 0;
int period = 1000;

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
  Serial.printf("%d;", counts);
  counts = 0;
  attachInterrupt(digitalPinToInterrupt(INTERRUPT_PIN), count, FALLING);
}
