package mqttclient

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/uber-go/zap"
)

// MqttClient contains the data needed to connect the client
type MqttClient struct {
	MqttServerHost string
	MqttServerPort int
	ConfigPath     string
	Config         *viper.Viper
	Logger         zap.Logger
	MqttClient     mqtt.Client
}

var client *MqttClient
var once sync.Once

// GetMqttClient creates the mqttclient and returns it
func GetMqttClient(configPath string, onConnectHandler mqtt.OnConnectHandler) *MqttClient {
	once.Do(func() {
		client = &MqttClient{
			ConfigPath: configPath,
			Config:     viper.New(),
		}
		client.configure()
		client.start(onConnectHandler)
	})
	return client
}

func (mc *MqttClient) configure() {
	mc.Logger = zap.NewJSON(zap.WarnLevel)

	mc.setConfigurationDefaults()
	mc.loadConfiguration()
	mc.configureClient()
}

func (mc *MqttClient) setConfigurationDefaults() {
	mc.Config.SetDefault("mqttserver.host", "localhost")
	mc.Config.SetDefault("mqttserver.port", 1883)
	mc.Config.SetDefault("mqttserver.user", "admin")
	mc.Config.SetDefault("mqttserver.pass", "admin")
}

func (mc *MqttClient) loadConfiguration() {
	mc.Config.SetConfigFile(mc.ConfigPath)
	mc.Config.SetEnvPrefix("mqttbridge")
	mc.Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	mc.Config.AutomaticEnv()

	if err := mc.Config.ReadInConfig(); err == nil {
		mc.Logger.Info("Loaded config file.", zap.String("configFile", mc.Config.ConfigFileUsed()))
	} else {
		panic(fmt.Sprintf("Could not load configuration file from: %s", mc.ConfigPath))
	}
}

func (mc *MqttClient) configureClient() {
	mc.MqttServerHost = mc.Config.GetString("mqttserver.host")
	mc.MqttServerPort = mc.Config.GetInt("mqttserver.port")
}

func (mc *MqttClient) start(onConnectHandler mqtt.OnConnectHandler) {
	logger.Logger.Debug("Initializing mqtt client")

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", mc.MqttServerHost, mc.MqttServerPort)).SetClientID("mqttbridge")
	opts.SetUsername(mc.Config.GetString("mqttserver.user"))
	opts.SetPassword(mc.Config.GetString("mqttserver.pass"))
	opts.SetKeepAlive(3 * time.Second)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetOnConnectHandler(onConnectHandler)
	mc.MqttClient = mqtt.NewClient(opts)

	c := mc.MqttClient

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Logger.Fatal(token.Error())
	}

	logger.Logger.Info(fmt.Sprintf("Successfully connected to mqtt server at %s:%d!",
		mc.MqttServerHost, mc.MqttServerPort))
}