// Copyright 2011 Aaron Jacobs. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A program that echos each byte read from the serial connection back
// across the connection, after incrementing it by one.

void setup() {
  Serial.begin(19200);
}

void loop() {
  if (Serial.available() > 0) {
    const uint8_t incoming_byte = Serial.read();
    Serial.write(incoming_byte + 1);
  }
}
