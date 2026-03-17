struct t_decodeCommand {
  String command;
  String value;
};

struct t_decodeRes {
  String command;
  bool isOk;
  String value;
};

struct t_state {
  bool beacon;
};

struct t_beacon_state {
  bool ledOn;
  int lastBlink;
};

t_state state;
t_beacon_state beacon_state;

void setup() {
  Serial.begin(115200);
  pinMode(LED_BUILTIN, OUTPUT);
}

void loop() {
  // Command handler
  if (Serial.available() > 0) {
    String encodedCommand = Serial.readString();
    encodedCommand.trim();

    if (encodedCommand == "_test") {
      String encodedCommandT = encodeCommand("foo", "bar");
      t_decodeCommand decodedCommandT = decodeCommand(encodedCommandT);
      String encodedResT = encodeRes("foo", true, "bar");
      t_decodeRes decodedResT = decodeRes(encodedResT);
      Serial.println("encode command gave: " + encodedCommandT);
      Serial.println("decode command gave: [command]: " + decodedCommandT.command + " [value]: " + decodedCommandT.value);
      Serial.println("encode res gave: " + encodedResT);
      Serial.println("decode res gave: [command]: " + decodedResT.command + " [isOk]: " + String(decodedResT.isOk) + " [value]: " + decodedResT.value);
      goto executor;
    }

    t_decodeCommand decodedCommand = decodeCommand(encodedCommand);

    if (decodedCommand.command == "" || decodedCommand.value == "") {
      String encodedRes = encodeRes("err", false, "invalid_command");
      Serial.println(encodedRes);
      goto executor;
    }

    if (decodedCommand.command == "ping") {
      String encodedRes = encodeRes(decodedCommand.command, true, "pong");
      Serial.println(encodedRes);
      goto executor;
    }

    if (decodedCommand.command == "led") {
      if (decodedCommand.value == "on") {
        digitalWrite(LED_BUILTIN, HIGH);
        String encodedRes = encodeRes(decodedCommand.command, true, "on");
        Serial.println(encodedRes);
        goto executor;
      }
      if (decodedCommand.value == "off") {
        digitalWrite(LED_BUILTIN, LOW);
        String encodedRes = encodeRes(decodedCommand.command, true, "off");
        Serial.println(encodedRes);
        goto executor;
      }
    }

    if (decodedCommand.command == "beacon") {
      if (decodedCommand.value == "on") {
        state.beacon = true;
        String encodedRes = encodeRes(decodedCommand.command, true, "on");
        Serial.println(encodedRes);
        goto executor;
      }
      if (decodedCommand.value == "off") {
        state.beacon = false;
        String encodedRes = encodeRes(decodedCommand.command, true, "off");
        Serial.println(encodedRes);
        goto executor;
      }
    }
  }

  // Exectutor handler
executor:
  if (state.beacon) {
    if (beacon_state.lastBlink == 0 || millis() - beacon_state.lastBlink >= 1000) {
      strobePattern();
    }
  }
};


void strobePattern() {
  for (int i = 0; i < 2; i++) {
    digitalWrite(LED_BUILTIN, HIGH);
    delay(50);
    digitalWrite(LED_BUILTIN, LOW);
    delay(50);
  };
  beacon_state.lastBlink = millis();
};

String encodeCommand(String command, String value) {
  return command + ":" + value;
}

t_decodeCommand decodeCommand(String encodedCommand) {
  t_decodeCommand decodedCommand;
  int ind0 = encodedCommand.indexOf(":");
  String command = encodedCommand.substring(0, ind0);
  command.trim();
  String value = encodedCommand.substring(ind0 + 1, encodedCommand.length());
  value.trim();
  decodedCommand.command = command;
  decodedCommand.value = value;
  return decodedCommand;
}

String encodeRes(String command, bool isOk, String value) {
  String encodedRes = command + ":";
  if (isOk) {
    encodedRes += "ok";
  } else {
    encodedRes += "not_ok";
  }
  encodedRes += ":" + value;
  return encodedRes;
}

t_decodeRes decodeRes(String encodedRes) {
  t_decodeRes decodedRes;
  int ind0 = encodedRes.indexOf(":");
  String command = encodedRes.substring(0, ind0);
  command.trim();
  int ind1 = encodedRes.indexOf(":", ind0 + 1);
  String status = encodedRes.substring(ind0 + 1, ind1 + 1);
  status.trim();
  String value = encodedRes.substring(ind1 + 1, encodedRes.length());
  value.trim();
  decodedRes.command = command;
  decodedRes.value = value;
  if (status == "ok") {
    decodedRes.isOk = true;
  } else {
    decodedRes.isOk = false;
  }
  return decodedRes;
}