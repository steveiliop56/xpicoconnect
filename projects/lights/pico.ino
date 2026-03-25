// Yes the code has comments, no AI did not generate them,
// I wrote them in order to make the flow easier to understand

// Structs holding the decoded command and response data
struct t_decodeCommand {
  String command;
  String value;
};

struct t_decodeResponse {
  String command;
  bool ok;
  String value;
};

// Main app state, holds values about the switches and sensors
struct t_switchStates {
  int beaconSwitchState = 0;
  int landSwitchState = 0;
  int taxiSwitchState = 0;
  int navSwitchState = 0;
  int strobeSwitchState = 0;
};

t_switchStates switchStates;

// So we don't spam the client if it's not ready to receive
bool isConnected = false;

int beaconSwitchPin = D0;
int landSwitchPin = D1;
int taxiSwitchPin = D2;
int navSwitchPin = D3;
int strobeSwitchPin = D4;

void setup() {
  Serial.begin(115200);

  pinMode(LED_BUILTIN, OUTPUT);
  pinMode(D0, INPUT_PULLDOWN);
  pinMode(D1, INPUT_PULLDOWN);
  pinMode(D2, INPUT_PULLDOWN);
  pinMode(D3, INPUT_PULLDOWN);
  pinMode(D4, INPUT_PULLDOWN);

  // Initial states
  switchStates.beaconSwitchState = digitalRead(D0);
  switchStates.landSwitchState = digitalRead(D1);
  switchStates.taxiSwitchState = digitalRead(D2);
  switchStates.navSwitchState = digitalRead(D3);
  switchStates.strobeSwitchState = digitalRead(D4);
}

void loop() {
  // Command handler
  if (Serial.available() > 0) {
    String encodedCommand = Serial.readStringUntil('\n');
    t_decodeCommand decodedCommand = decodeCommand(encodedCommand);

    // Validate command
    if (decodedCommand.command == "" || decodedCommand.value == "") {
      String encodedResponse = encodeResponse("err", false, "invalid_command");
      Serial.println(encodedResponse);
    }

    // Simple ping, mostly for testing
    if (decodedCommand.command == "ping") {
      String encodedResponse = encodeResponse(decodedCommand.command, true, "pong");
      Serial.println(encodedResponse);
    }

    // FDX means full-duplex and it confirms that the client is ready to receive commands and responses.
    // It's a handshake command that should be sent at the beginning of the connection.
    if (decodedCommand.command == "fdx") {
      String encodedResponse = encodeResponse(decodedCommand.command, true, decodedCommand.value);
      Serial.println(encodedResponse);

      // Give some time to the client to process the response before we start sending commands
      delay(50);

      String encodedCommand = encodeCommand("fdx", decodedCommand.value);
      Serial.println(encodedCommand);

      Serial.setTimeout(2000);
      String clientEncodedResponse = Serial.readStringUntil('\n');
      t_decodeResponse decodedResponse = decodeResponse(clientEncodedResponse);

      if (!decodedResponse.ok || decodedResponse.command != "fdx" || decodedResponse.value != "bar") {
        String encodedError = encodeResponse("err", false, "fdx_failed");
        Serial.println(encodedError);
      }

      // Since we confirmed full-duplex, we are connected
      isConnected = true;
      digitalWrite(LED_BUILTIN, HIGH);
    }

    // End command "disconnects" the client
    if (decodedCommand.command == "end") {
      digitalWrite(LED_BUILTIN, LOW);
      isConnected = false;

      String encodedResponse = encodeResponse(decodedCommand.command, true, "ok");
      Serial.println(encodedResponse);
    }


    // Led just shows the user that we are connected
    if (decodedCommand.command == "led") {
      if (decodedCommand.value == "on") {
        digitalWrite(LED_BUILTIN, HIGH);
        String encodedResponse = encodeResponse(decodedCommand.command, true, "on");
        Serial.println(encodedResponse);
      }
      if (decodedCommand.value == "off") {
        digitalWrite(LED_BUILTIN, LOW);
        String encodedResponse = encodeResponse(decodedCommand.command, true, "off");
        Serial.println(encodedResponse);
      }
    }
  }

  // After command checking, we can see what we need to send back
  if (!isConnected) {
    return;
  };

  handleSwitchState(beaconSwitchPin, switchStates.beaconSwitchState, "beacon_switch");
  handleSwitchState(landSwitchPin, switchStates.landSwitchState, "land_switch");
  handleSwitchState(taxiSwitchPin, switchStates.taxiSwitchState, "taxi_switch");
  handleSwitchState(navSwitchPin, switchStates.navSwitchState, "nav_switch");
  handleSwitchState(strobeSwitchPin, switchStates.strobeSwitchState, "strobe_switch");

  delay(50);
};

void handleSwitchState(int pin, int &switchState, String switchName) {
  int switchValue = digitalRead(pin);
  if (switchValue != switchState) {
    switchState = switchValue;
    String encodedCommand = encodeCommand(switchName, String(switchValue));
    Serial.println(encodedCommand);
  };
};

String encodeCommand(String command, String value) {
  return command + ":" + value;
};

t_decodeCommand decodeCommand(String encodedCommand) {
  t_decodeCommand decodedCommand;
  encodedCommand.trim();

  int ind0 = encodedCommand.indexOf(":");
  String command = encodedCommand.substring(0, ind0);
  command.trim();

  String value = encodedCommand.substring(ind0 + 1, encodedCommand.length());
  value.trim();

  decodedCommand.command = command;
  decodedCommand.value = value;

  return decodedCommand;
};

String encodeResponse(String command, bool ok, String value) {
  String encodedResponse = command + ":";
  if (ok) {
    encodedResponse += "ok";
  } else {
    encodedResponse += "not_ok";
  }
  encodedResponse += ":" + value;
  return encodedResponse;
};

t_decodeResponse decodeResponse(String encodedResponse) {
  t_decodeResponse decodedResponse;
  encodedResponse.trim();

  int ind0 = encodedResponse.indexOf(":");
  String command = encodedResponse.substring(0, ind0);
  command.trim();

  int ind1 = encodedResponse.indexOf(":", ind0 + 1);
  String status = encodedResponse.substring(ind0 + 1, ind1);
  ;
  status.trim();

  String value = encodedResponse.substring(ind1 + 1, encodedResponse.length());
  value.trim();

  decodedResponse.command = command;
  decodedResponse.value = value;

  if (status == "ok") {
    decodedResponse.ok = true;
  } else {
    decodedResponse.ok = false;
  }

  return decodedResponse;
};
