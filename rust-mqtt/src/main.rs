use rumqttc::MqttOptions;

fn main() {
    let mut mqttOptions = MqttOptions::new("logger", "ip", 8333);
}
