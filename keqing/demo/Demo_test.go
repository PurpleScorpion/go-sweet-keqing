package demo

import (
	"errors"
	"fmt"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"testing"
	"time"
)

type animal struct {
	Jinmao dog
	Dogs   []dog
}

type dog struct {
	Name string
	Age  int
}

func TestDemo31(t *testing.T) {
	filePath := "D:\\img\\哈哈.jpg"
	savePath := "D:\\imgs\\haha\\哈哈.jpg"
	err := keqing.CopyFile(filePath, savePath)
	if err != nil {
		fmt.Println(err)
	}
}
func TestDemo30(t *testing.T) {

	random := keqing.WinningRate(0.2)
	random2 := keqing.WinningRate(0.02)
	random3 := keqing.WinningRate(80)
	fmt.Println(random, random2, random3)
}
func TestDemo29(t *testing.T) {

	random := keqing.Random()
	for i := 0; i < 10; i++ {
		random1 := keqing.RandomNextInt(100)
		fmt.Println(random1)
	}
	random2 := keqing.RandomNextInt(100, 1000)
	random3 := keqing.RandomId()
	fmt.Println(random, random2, random3)
}
func TestDemo28(t *testing.T) {

	str := "ems-dev"
	flag := keqing.StartWith(str, "ems-1")
	fmt.Println(flag)
}
func TestDemo27(t *testing.T) {

	//days := keqing.GetDaysInMonth(2025, 2)
	//days := keqing.GetDaysInMonth(time.Now())
	now := time.Date(2025, 2, 1, 0, 0, 0, 0, time.Local)
	days := keqing.GetDaysInMonth(now)
	fmt.Println(days)

}

func TestDemo26(t *testing.T) {
	str := "elec_gen_coef_gj_kWh\telec_gen_co2_coef\telec_day_coef\telec_day_co2_coef\telec_night_coef\telec_night_co2_coef\telec2_gen_coef\telec2_gen_co2_coef\telec2_day_coef\telec2_day_co2_coef\telec2_night_coef\telec2_night_co2_coef\tcity_gas_coef\tcity_gas_co2_coef\tlp_gas_coef\tlp_gas_co2_coef\toil_a_heavy_coef\toil_a_heavy_co2_coef\toil_light_coef\toil_light_co2_coef\toil_kero_coef\toil_kero_co2_coef\tdistrict_energy_coef\tdistrict_energy_co2_coef\tother1_energy_coef\tother1_co2_coef\tother2_energy_coef\tother2_co2_coef\tsolar_self_consumption_coef\tsolar_self_consumption_co2\tsolar_sell_coef\tsolar_sell_co2\tsolar_total_gen_coef\tsolar_total_gen_co2\tcogeneration_elec_coef\tcogeneration_elec_co2\trenewable_self_consumption_coef\trenewable_self_consumption_co2\trenewable_sell_coef\trenewable_sell_co2\tdistrict_heating_coef\tdistrict_heating_co2\tother1_energy_coef_kWh\tother1_co2_coef_kWh\tother2_energy_coef_kWh\tother2_co2_coef_kWh"
	arr := keqing.Split(str, "\t")

	for i := 0; i < len(arr); i++ {
		char := toExcelChar(i+6) + "21"
		fmt.Println(char + "\t\t" + keqing.Snake2BigCamel(arr[i]) + "\t\t\t\t" + "get" + keqing.Snake2BigCamel(arr[i]))
	}

}

func TestDemo25(t *testing.T) {
	str := "other_elec_gen_gj\tother_elec_day_gj\tother_elec_night_gj\tother_city_gas_gj\tother_lp_gas_gj"
	arr := keqing.Split(str, "\t")

	for i := 0; i < len(arr); i++ {
		char := toExcelChar(i+6) + "21"
		fmt.Println(arr[i] + "\t\t\t\t" + char + "\t\t\t\t" + keqing.Snake2BigCamel(arr[i]) + "\t\t\t\t" + "get" + keqing.Snake2BigCamel(arr[i]))
	}

}

func toExcelChar(n int) string {
	var s string
	for n > 0 {
		n--
		s = string(rune('A'+n%26)) + s
		n /= 26
	}
	return s
}

func TestDemo24(t *testing.T) {
	str := "heat_elec_gen_gj\theat_elec_day_gj\theat_elec_night_gj\theat_city_gas_gj\theat_lp_gas_gj\theat_a_heavy_oil_gj\theat_light_oil_gj\theat_kero_gj\theat_hvac_gj\theat_other1_gj\tac_elec_gen_gj\tac_elec_day_gj\tac_elec_night_gj\tpump_elec_gen_gj\tpump_elec_day_gj\tpump_elec_night_gj\tac_std_energy_gj\tac_elec_total\tac_reduction_gj\tac_reduction_rate\tac_bei\tcogeneration_city_gas_gj\tcogeneration_lp_gas_gj\tcogeneration_a_heavy_oil_gj\tcogeneration_light_oil_gj\tcogeneration_kero_gj\tvent_elec_gen_gj\tvent_elec_day_gj\tvent_elec_night_gj\tvent_std_energy_gj\tvent_elec_total\tvent_reduction_gj\tvent_reduction_rate\tvent_bei\tlight_elec_gen_gj\tlight_elec_day_gj\tlight_elec_night_gj\tlight_std_energy_gj\tlight_elec_total\tlight_reduction_gj\tlight_reduction_rate\tlight_bei\thotwater_elec_gen_gj\thotwater_elec_day_gj\thotwater_elec_night_gj\thotwater_city_gas_gj\thotwater_lp_gas_gj\thotwater_a_heavy_oil_gj\thotwater_light_oil_gj\thotwater_kero_gj\thotwater_other2_gj\thotwater_std_energy_gj\thotwater_elec_total\thotwater_reduction_gj\thotwater_reduction_rate\thotwater_bei\televator_elec_gen_gj\televator_elec_day_gj\televator_elec_night_gj\televator_std_energy_gj\televator_elec_total\televator_reduction_gj\televator_reduction_rate\televator_bei\tother_elec_gen_gj\tother_elec_day_gj\tother_elec_night_gj\tother_city_gas_gj\tother_lp_gas_gj\tother_std_energy_gj\tother_elec_total\tother_reduction_gj\tsolar_self_consumption1_gj\tcogeneration_elec_gen3_gj\trenewable_self_consumption1_gj\trenewable_total_generation_gj\tcogeneration_std_energy_gj\tcogeneration_reduction_gj\tcogeneration_reduction_rate\tstd_energy_gj\telec_total\treduction_gj\treduction_rate\tbei\tmonth_reduction_rate"
	arr := keqing.Split(str, "\t")

	for i := 0; i < len(arr); i++ {
		s := keqing.ToLowerCase(arr[i])
		fmt.Println(s + ",")
	}

}

func TestDemo22(t *testing.T) {
	var str = `heat_elec_gen_gj,
heat_elec_day_gj,
heat_elec_night_gj,
heat_city_gas_gj,
heat_lp_gas_gj,
heat_a_heavy_oil_gj,
heat_light_oil_gj,
heat_kero_gj,
heat_hvac_gj,
heat_other1_gj,
ac_elec_gen_gj,
ac_elec_day_gj,
ac_elec_night_gj,
pump_elec_gen_gj,
pump_elec_day_gj,
pump_elec_night_gj,
ac_std_energy_gj,
ac_elec_total,
ac_reduction_gj,
ac_reduction_rate,
ac_bei,
cogeneration_city_gas_gj,
cogeneration_lp_gas_gj,
cogeneration_a_heavy_oil_gj,
cogeneration_light_oil_gj,
cogeneration_kero_gj,
vent_elec_gen_gj,
vent_elec_day_gj,
vent_elec_night_gj,
vent_std_energy_gj,
vent_elec_total,
vent_reduction_gj,
vent_reduction_rate,
vent_bei,
light_elec_gen_gj,
light_elec_day_gj,
light_elec_night_gj,
light_std_energy_gj,
light_elec_total,
light_reduction_gj,
light_reduction_rate,
light_bei,
hotwater_elec_gen_gj,
hotwater_elec_day_gj,
hotwater_elec_night_gj,
hotwater_city_gas_gj,
hotwater_lp_gas_gj,
hotwater_a_heavy_oil_gj,
hotwater_light_oil_gj,
hotwater_kero_gj,
hotwater_other2_gj,
hotwater_std_energy_gj,
hotwater_elec_total,
hotwater_reduction_gj,
hotwater_reduction_rate,
hotwater_bei,
elevator_elec_gen_gj,
elevator_elec_day_gj,
elevator_elec_night_gj,
elevator_std_energy_gj,
elevator_elec_total,
elevator_reduction_gj,
elevator_reduction_rate,
elevator_bei,
other_elec_gen_gj,
other_elec_day_gj,
other_elec_night_gj,
other_city_gas_gj,
other_lp_gas_gj,
other_std_energy_gj,
other_elec_total,
other_reduction_gj,
solar_self_consumption1_gj,
cogeneration_elec_gen3_gj,
renewable_self_consumption1_gj,
renewable_total_generation_gj,
cogeneration_std_energy_gj,
cogeneration_reduction_gj,
cogeneration_reduction_rate,
std_energy_gj,
elec_total,
reduction_gj,
reduction_rate,
bei,
month_reduction_rate`
	arr := keqing.Split(str, ",")
	for i := 0; i < len(arr); i++ {
		fmt.Println(fmt.Sprintf("pojo.%s,", keqing.Trim(keqing.Snake2BigCamel(arr[i]))))
	}
	fmt.Println("============================")
	for i := 0; i < len(arr); i++ {
		fmt.Print("?, ")
	}
	fmt.Println("============================")
	for i := 0; i < len(arr); i++ {
		fmt.Println(fmt.Sprintf("%s = ?,", keqing.Trim(arr[i])))
	}
}

func TestDemo21(t *testing.T) {
	var publicKeyStrEv = `-----BEGIN PUBLIC KEY-----
公钥
-----END PUBLIC KEY-----`

	var privateKeyStrEv = `-----BEGIN RSA PRIVATE KEY-----
私钥
-----END RSA PRIVATE KEY-----`
	rsaLoadPublicKey := keqing.RsaLoadPublicKey(publicKeyStrEv)
	rsaLoadPrivateKey := keqing.RsaLoadPrivateKey(privateKeyStrEv)

	token := keqing.RsaEncrypt4PKCS1(rsaLoadPublicKey, "hahah")
	fmt.Println(token)
	data := keqing.RsaDecrypt4PKCS1(rsaLoadPrivateKey, token)
	fmt.Println(data)
}

func TestDemo20(t *testing.T) {
	var str1 = "2024-08-30 10:48:01"
	var str2 = "2024-07-19T21:30:04.175576Z"
	fmt.Println("Local => Local", keqing.ParseLocal(str1))
	fmt.Println("Local => UTC", keqing.Local2UTC(str1))
	fmt.Println("UTC => UTC", keqing.ParseUTC(str2))
	fmt.Println("UTC => Local", keqing.UTC2Local(str2))

	fmt.Println("Parse -------------------------------- ")
	fmt.Println("Local ", keqing.ParseDate(str1, keqing.DEFAULT_LOCAL_FORMAT, keqing.LOCAL))
	fmt.Println("UTC ", keqing.ParseDate(str2, keqing.DEFAULT_UTC_FORMAT, keqing.UTC))
}

func TestDemo19(t *testing.T) {
	var dongwu animal
	dongwu.Jinmao.Name = "金毛"
	dongwu.Jinmao.Age = 18
	var dogs []dog
	dogs = append(dogs, dog{Name: "二哈", Age: 5})
	dogs = append(dogs, dog{Name: "边牧", Age: 6})
	dogs = append(dogs, dog{Name: "泰迪", Age: 7})
	dongwu.Dogs = dogs
	fmt.Println(keqing.ToString(dongwu))
}

func TestDemo18(t *testing.T) {
	mymap := make(map[string]interface{})
	mymap["name"] = "keqing"
	mymap["age"] = 13
	mymap["hasGo"] = true
	mymap["list"] = []string{"banana", "apple", "orange", "grape"}
	mymap["dog"] = dog{
		Name: "二哈",
		Age:  5,
	}
	objList := make([]interface{}, 0)
	objList = append(objList, "张三")
	objList = append(objList, "李四")
	objList = append(objList, 100)
	objList = append(objList, true)
	objList = append(objList, dog{Name: "狗子", Age: 6})
	objList = append(objList, []string{"阿狸", "亚索", "提莫"})
	mymap["objList"] = objList

	fmt.Println(keqing.ToString(mymap))

}

func TestDemo17(t *testing.T) {
	keqing.RsaLoadKey("D:\\img\\publicKey.pem", "D:\\img\\privateKey.pem", keqing.RSA_KEY_FILE_TYPE)
	str := keqing.RsaEncrypt("剑光如我,斩尽牛杂")
	fmt.Println(str)
	fmt.Println("-------------------")
	data := keqing.RsaDecrypt(str)
	fmt.Println(data)
}
func TestDemo16(t *testing.T) {
	keqing.RsaGenerateKey("D:\\img")
}

func TestDemo15(t *testing.T) {
	publicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuOP6JK8tT0WdIvMbK60T
x/BrbuZN/tJn+HVXwh9xdDAyygMhpI7yfX1NFzkAtoBb9dCEbmvAtImSlq+grDMQ
8qse/KI0W0viQLbIXwUdIaSx/qKYjUprYgBwlMN8woV99IB+SwZ7T5SHMb1LK6xM
9ybaeiNxOPoeGy4S6YvEpC1tjtKKq8kpBlaN3psx5H2jeR//Jj/h4V/7trzWO5Em
LgcC9jx4OazrAR85jOKgLo9L0KmKe6S4icQ8It7LN17BciqUKMrYQS0kj6R+HtoI
3f79wNBLdTna3DYsF3YDLBrlPV8GB0PSTUTqNSoX42/qlZVuH8wz8TGpQijsRc/v
3wIDAQAB
-----END PUBLIC KEY-----`
	privateKey := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC44/okry1PRZ0i
8xsrrRPH8Gtu5k3+0mf4dVfCH3F0MDLKAyGkjvJ9fU0XOQC2gFv10IRua8C0iZKW
r6CsMxDyqx78ojRbS+JAtshfBR0hpLH+opiNSmtiAHCUw3zChX30gH5LBntPlIcx
vUsrrEz3Jtp6I3E4+h4bLhLpi8SkLW2O0oqrySkGVo3emzHkfaN5H/8mP+HhX/u2
vNY7kSYuBwL2PHg5rOsBHzmM4qAuj0vQqYp7pLiJxDwi3ss3XsFyKpQoythBLSSP
pH4e2gjd/v3A0Et1OdrcNiwXdgMsGuU9XwYHQ9JNROo1Khfjb+qVlW4fzDPxMalC
KOxFz+/fAgMBAAECggEAHYDzt8rkdhPrwVn96fhSgcNRwX6qz5EP2kwPVwDhf+L5
F9dsFPBirbfDB4OnI3hUNGOz3lL/i0+wvq8D+raja7X22eWgaTkwv5brXo5YWbgI
V1Pm+BT6Ecd0L6kKTZgzw0KF5L8CCm7vK/bC+hMirQXcM0VYmfj/uOKfTflpxbDc
OPtzT99lVhjJg2e7GTnwpIXERI8SByjc+L1pQrN4FsoOl4VTg8o+QeLcKHghFaZ1
jnQb9B2EQ1TrsGnww1tNrNOvV9k/UBXCH+A3s6OCaHfNLoREATMqrtX11a2qzMZz
7VZtWC7uqI/FjZZEDluXYlJHqe9VY4c7cVqMIuNSMQKBgQDz72qDrdufi5z61h89
vWg0giIC6gkyy6OhwDYbv4ydI/9m3vexGbkX33YLmSxxKmgiYT6UD+LvjiHw6Waa
0iFX7L0FF1KgPRxhfs1g37YPWD1e/2tpT17VKLqjMAX4xOwZ3HfVvZBI/0g6hjzZ
sBB3Vtq0iyU+yXFZ5u2pgtqp7QKBgQDCCPec900fNF6e08Z0fxWVYqPiOw4I3iWb
emi3pGPWMLPoOwKDPfguSaHJKhi/Gq9Ipd30T2xkvuW7jMFzMOcGhoVMmFVuLfJz
ttj8J5x1YtgmNrXYwtWXQJypqmZJTYwXOBzUYRiplIZNvA1FIy0tSVyH1Zg85/43
vpElJ+wXewKBgHAOuaWIBm4CWri4CF36VpZYeXtRO6yD88VoYPLaSaQeV0NQhgRr
RqX612V4lfveeTvh5DdsHNnjNyBOd/4DLaIQdLyT/Db0G8eF0p7/5ciixn6PYy5b
cbsGHMa+Vt/yxmsS5lHf5RpDe1C3PdjakpXf5lQt34w6ScH83YyTOhP5AoGBAK5C
s28LQv4lYF0wQOlbQR0aq6h/9QjNyeSquOVFBEzXDJwicw0/WGbpxh0Oa48l/go2
vPGvat/H+jbIIOy9HJ7lrU2u+fqr1TVLH/DF+mQKU6luNT7pLD5cztYprRdkR86K
nIm4chfKxhuGKjzPbMFhQ3LSx5jbmZqi0WQXSJeFAoGARorE8VlMqDbgtSP/TApp
HAr/UuVdmf8Ag+ZUpa2o33E39L9yLnj7L60GZKVJ87weYj/CHatiCIzpStNsGkoU
81iOD/cUhLl7K9B0WeT7ECfO4Vk6xg6n0SH6YykEyudAmlrtVfxZ1EclQ2OKhxfN
Ia8OCVyYe7VFKVYFkh4PBhQ=
-----END PRIVATE KEY-----`
	str := keqing.RsaEncryptCustom(publicKey, "剑光如我,斩尽牛杂")
	fmt.Println(str)
	fmt.Println("-------------------")
	data := keqing.RsaDecryptCustom(privateKey, str)
	fmt.Println(data)
}

func demo14() {
	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 15, 20, 215, 30, 35, 27, 40, 55, 77}

	arr := keqing.Difference(intSlice1, intSlice2)
	fmt.Println(keqing.ToString(arr))

	stringSlice := []string{"banana", "apple1", "orange1", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	arr2 := keqing.Difference(stringSlice, expectedStringSlice)
	fmt.Println(keqing.ToString(arr2))
}

func demo13() {
	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 15, 20, 215, 30, 35, 27, 40, 55, 77}

	arr := keqing.Union(intSlice1, intSlice2)
	fmt.Println(keqing.ToString(arr))

	stringSlice := []string{"banana", "apple1", "orange1", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	arr2 := keqing.Union(stringSlice, expectedStringSlice)
	fmt.Println(keqing.ToString(arr2))
}

func demo12() {
	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 15, 20, 215, 30, 35, 27, 40, 55, 77}

	arr := keqing.Intersection(intSlice1, intSlice2)
	fmt.Println(keqing.ToString(arr))

	stringSlice := []string{"banana", "apple1", "orange1", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	arr2 := keqing.Intersection(stringSlice, expectedStringSlice)
	fmt.Println(keqing.ToString(arr2))
}

func demo11() {

	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	expectedIntSlice1 := []int32{9, 12, 27, 31, 45, 55, 62, 77, 85, 96}
	expectedIntSlice2 := []int64{9, 12, 27, 31, 45, 55, 62, 77, 85, 96}

	keqing.SortArrayDesc(&intSlice1)
	keqing.SortArrayDesc(&intSlice2)
	keqing.SortArrayDesc(&expectedIntSlice1)
	keqing.SortArrayDesc(&expectedIntSlice2)
	fmt.Println(keqing.ToString(intSlice1))
	fmt.Println(keqing.ToString(intSlice2))
	fmt.Println(keqing.ToString(expectedIntSlice1))
	fmt.Println(keqing.ToString(expectedIntSlice2))

	stringSlice := []string{"banana", "apple", "orange", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	keqing.SortArrayDesc(&stringSlice)
	keqing.SortArrayDesc(&expectedStringSlice)
	fmt.Println(keqing.ToString(stringSlice))
	fmt.Println(keqing.ToString(expectedStringSlice))

}

func demo10() {

	intSlice1 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	intSlice2 := []int{12, 85, 96, 45, 62, 9, 27, 31, 55, 77}
	expectedIntSlice1 := []int32{9, 12, 27, 31, 45, 55, 62, 77, 85, 96}
	expectedIntSlice2 := []int64{9, 12, 27, 31, 45, 55, 62, 77, 85, 96}

	keqing.SortArray(&intSlice1)
	keqing.SortArray(&intSlice2)
	keqing.SortArray(&expectedIntSlice1)
	keqing.SortArray(&expectedIntSlice2)
	fmt.Println(keqing.ToString(intSlice1))
	fmt.Println(keqing.ToString(intSlice2))
	fmt.Println(keqing.ToString(expectedIntSlice1))
	fmt.Println(keqing.ToString(expectedIntSlice2))

	stringSlice := []string{"banana", "apple", "orange", "grape"}
	expectedStringSlice := []string{"apple", "banana", "grape", "orange"}
	keqing.SortArray(&stringSlice)
	keqing.SortArray(&expectedStringSlice)
	fmt.Println(keqing.ToString(stringSlice))
	fmt.Println(keqing.ToString(expectedStringSlice))

}

func demo9() {
	fmt.Println("NowDate:  ", keqing.NowDate())
	fmt.Println("NowUTCDate:  ", keqing.NowUTCDate())
	fmt.Println("FormatNowDate:  ", keqing.NowDateStr())
	fmt.Println("FormatNowUTCDate:  ", keqing.NowUTCDateStr())
	fmt.Println("CurrentTimeMillis:  ", keqing.CurrentTimeMillis())

	fmt.Println("NowLocalDate:  ", keqing.ToString(keqing.DateInfo(keqing.NowDate())))
	fmt.Println("NowUTCDate:  ", keqing.ToString(keqing.DateInfo(keqing.NowUTCDate())))
}

func demo8() {

	var arr = make([]interface{}, 0)
	arr = append(arr, "a")
	arr = append(arr, "2")
	arr = append(arr, 3)
	arr = append(arr, true)
	arr = append(arr, false)

	fmt.Println(keqing.ToString(arr))

	//var obj = dog{
	//	Age:  1,
	//	Name: "dog",
	//}
	//fmt.Println(keqing.ToString(obj))
}

func demo7() {
	str := "  123456  "
	fmt.Println("AAA" + keqing.Trim(str) + "BBB")

	fmt.Println("TrimSuffix: ", keqing.TrimSuffix("123456789", "789"))
	fmt.Println("TrimSuffix: ", keqing.TrimSuffix("123456789", "9"))
	fmt.Println("TrimSuffix: ", keqing.TrimSuffix("123456789", "asd"))

	fmt.Println("Split: ", keqing.Split("123,456,789", ","))
	fmt.Println("ToFirstUpperCase: ", keqing.ToFirstUpperCase("a123456"))
	fmt.Println("ToFirstUpperCase: ", keqing.ToFirstUpperCase("123456"))
	fmt.Println("ToFirstLowerCase: ", keqing.ToFirstLowerCase("B123456"))
	fmt.Println("ToFirstLowerCase: ", keqing.ToFirstLowerCase("123456"))
	fmt.Println("ToUpperCase: ", keqing.ToUpperCase("123aWSsda456"))
	fmt.Println("ToLowerCase: ", keqing.ToLowerCase("12ADWsaSA3456"))

	fmt.Println("ToStringArray: ", keqing.ToStringArray("剑光如网,斩尽芜杂"))

}

func demo6() {
	for i := 0; i < 20; i++ {
		fmt.Println("UUID: ", keqing.UUID())
	}

}

func demo5() {
	a := "simple_test"
	b := "http_request"

	c := "JsonResponse"
	d := "xmlHttpRequest"

	fmt.Println("下划线转驼峰: ", keqing.Snake2Camel(a))
	fmt.Println("下划线转大驼峰: ", keqing.Snake2BigCamel(b))

	fmt.Println("大驼峰转下划线", keqing.Camel2Snake(c))
	fmt.Println("小驼峰转下划线", keqing.Camel2Snake(d))

}

func demo4() {
	fmt.Println("MD5: ", keqing.MD5("123456"))
	fmt.Println("MD5Salt: ", keqing.MD5Salt("123456", "asd"))
}

func demo3() {
	var a int = 10
	var b int8 = 20
	var c int16 = 30
	var d float32 = 40.1234
	var e float64 = 50.123456789

	fmt.Println("int-string", keqing.Num2Str(a))
	fmt.Println("int8-string", keqing.Num2Str(b))
	fmt.Println("int16-string", keqing.Num2Str(c))
	fmt.Println("float32-string", keqing.Num2Str(d))
	fmt.Println("float64-string", keqing.Num2Str(e))

	fmt.Println("float32-string-2", keqing.Float2Str(d, 2))
	fmt.Println("float64-string-4", keqing.Float2Str(e, 4))

}

// emptyData
func demo1() {
	var a int = 1
	var a1 int = 0
	fmt.Println("int-非空....", keqing.IsEmpty(a))
	fmt.Println("int-空....", keqing.IsEmpty(a1))

	var b int8 = 1
	var b1 int8 = 0
	fmt.Println("int8-非空....", keqing.IsEmpty(b))
	fmt.Println("int8-空....", keqing.IsEmpty(b1))

	var c int16 = 1
	var c1 int16 = 0
	fmt.Println("int8-非空....", keqing.IsEmpty(c))
	fmt.Println("int8-空....", keqing.IsEmpty(c1))

	d := make(map[string]int)
	d["one"] = 1
	d["two"] = 2
	d1 := make(map[string]int)
	fmt.Println("map-非空....", keqing.IsEmpty(d))
	fmt.Println("map-空....", keqing.IsEmpty(d1))

	var e []string
	var e1 []string
	e = append(e, "1")
	e = append(e, "2")
	fmt.Println("切片-非空....", keqing.IsEmpty(e))
	fmt.Println("切片-空....", keqing.IsEmpty(e1))

	var f = "a"
	var f1 = ""
	fmt.Println("string-非空....", keqing.IsEmpty(f))
	fmt.Println("string-空....", keqing.IsEmpty(f1))

	var err = errors.New("aaaaa")

	fmt.Println("nil值的变量-非空....", keqing.IsEmpty(err))
	fmt.Println("nil值的变量-空....", keqing.IsEmpty(getNilError()))

	//var obj = dog{
	//	Age:  1,
	//	Name: "dog",
	//}
	//fmt.Println("结构体....", keqing.IsEmpty(obj))
}

func demo2() {

	arr := make([]string, 0)
	arr = append(arr, "1")
	arr = append(arr, "2")
	arr = append(arr, "3")

	a := keqing.ArrayContains(arr, "2")
	a1 := keqing.ArrayContains(arr, "4")

	fmt.Println("包含: ", a)
	fmt.Println("不包含: ", a1)
}
func getNilError() error {
	return nil
}
