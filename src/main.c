#include <stdio.h>

#include <mosquitto.h>

void on_connect(struct mosquitto *mosq, void *obj, int st) { printf("connected"); }

void on_message(struct mosquitto *mosq, void *obj, const struct mosquitto_message *msg) { printf("message"); }

int main(void) {
    int st;
    mosquitto_lib_init();

    struct mosquitto *mosq;
    mosq = mosquitto_new("test", true, NULL);

    mosquitto_connect_callback_set(mosq, on_connect);
    mosquitto_message_callback_set(mosq, on_message);

    mosquitto_tls_set(mosq, "blert.pem", NULL, NULL, NULL, NULL);
    mosquitto_username_pw_set(mosq, "USER", "PASS");

    st = mosquitto_connect(mosq, "IP", 8883, 10);
    if (st) {
        printf("Could not connect with code: %d\n", st);
        return -1;
    }

    mosquitto_loop_start(mosq);
    printf("Press Enter to quit...\n");
    getchar();
    mosquitto_loop_stop(mosq, true);

    mosquitto_disconnect(mosq);
    mosquitto_destroy(mosq);
    mosquitto_lib_cleanup();

    return 0;
}
