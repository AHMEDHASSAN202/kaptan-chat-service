package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	. "github.com/ahmetb/go-linq/v3"
	"io"
	"os"
	"samm/internal/module/common/dto"
	"samm/pkg/logger"
	"strings"
)

const Cities = `[
    {
        "_id": "6666c1f317bf1d5b07e6fe85",
        "name": {
            "ar": "أبها",
            "en": "ABHA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3a731cb56dfe6fe86",
        "name": {
            "ar": "بقيق",
            "en": "ABQAIQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3fdcc89f06de6fe87",
        "name": {
            "ar": "ابو عريش",
            "en": "ABU ARISH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f379429c5be7e6fe88",
        "name": {
            "ar": "ابو حدرية",
            "en": "ABU HADRIYAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31819bd59dfe6fe89",
        "name": {
            "ar": "عفيف",
            "en": "AFIF"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f314425cb75de6fe8a",
        "name": {
            "ar": "الأفلاج",
            "en": "AFLAJ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3a28a2e4412e6fe8b",
        "name": {
            "ar": "احد مسرة",
            "en": "AHAD MASARAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3fe6ba8739fe6fe8c",
        "name": {
            "ar": "عهد رفيدة",
            "en": "AHAD ROFAIDAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f399320ccc1ee6fe8d",
        "name": {
            "ar": "عين دار",
            "en": "AIN DAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3115d73e09de6fe8e",
        "name": {
            "ar": "اجفار",
            "en": "AJFAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e1475a4d87e6fe8f",
        "name": {
            "ar": "أجياد",
            "en": "AJYAD"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c8080100fbe6fe90",
        "name": {
            "ar": "العنبرية",
            "en": "AL ANBARIYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f53080794de6fe91",
        "name": {
            "ar": "العريسة",
            "en": "AL ARISA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ad772ef55de6fe92",
        "name": {
            "ar": "الآسية",
            "en": "AL ASIYAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f34edc14f382e6fe93",
        "name": {
            "ar": "العوايلة",
            "en": "AL ATAWILAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f35798137f3de6fe94",
        "name": {
            "ar": "العوالي",
            "en": "AL AWALY"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ab9ef07866e6fe95",
        "name": {
            "ar": "العيس",
            "en": "AL AYSS"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f382f29f1f11e6fe96",
        "name": {
            "ar": "الباحة",
            "en": "AL BAHA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3870d632ec0e6fe97",
        "name": {
            "ar": "البحر",
            "en": "AL BAHAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f327ded015fbe6fe98",
        "name": {
            "ar": "البشير",
            "en": "AL BASHAIR"
        }
    },
    {
        "_id": "6666c1f3a6adb39de7e6fe99",
        "name": {
            "ar": "البسيتة",
            "en": "AL BUSAITA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38afee52c4ce6fe9a",
        "name": {
            "ar": "القزاز البلد",
            "en": "AL GAZAZ AL BALAD"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f375a48a11bbe6fe9b",
        "name": {
            "ar": "الجويزة",
            "en": "AL GEWAIZAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38e160b830fe6fe9c",
        "name": {
            "ar": "الغاط",
            "en": "AL GHAT"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b3efabf399e6fe9d",
        "name": {
            "ar": "غزة",
            "en": "AL GHAZAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f37077520a78e6fe9e",
        "name": {
            "ar": "الغزالة",
            "en": "AL GHAZALAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f77a21071ce6fe9f",
        "name": {
            "ar": "القوز",
            "en": "AL GOZ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f329262c6127e6fea0",
        "name": {
            "ar": "الحيط",
            "en": "AL HAET"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b4188311fde6fea1",
        "name": {
            "ar": "الحجون",
            "en": "AL HAJOUN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f371f17f9a5ce6fea2",
        "name": {
            "ar": "الحريق",
            "en": "AL HAREEQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ebf57ecc72e6fea3",
        "name": {
            "ar": "الهرجة",
            "en": "AL HARJAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3469c43637ee6fea4",
        "name": {
            "ar": "الحسا",
            "en": "AL HASA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f321c36300bfe6fea5",
        "name": {
            "ar": "الحرة الشرقية",
            "en": "AL HURRA AL SHARQYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c371e55f6ae6fea6",
        "name": {
            "ar": "الجوف",
            "en": "AL JOUF"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f37e097e06fbe6fea7",
        "name": {
            "ar": "الجهيمة",
            "en": "AL JUHAIMA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f335cbd984f7e6fea8",
        "name": {
            "ar": "الجموم",
            "en": "AL JUMUM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f30fbab30c01e6fea9",
        "name": {
            "ar": "الخرج",
            "en": "AL KHARJ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3fecf25ffede6feaa",
        "name": {
            "ar": "الخبر",
            "en": "AL KHOBAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e3b27eed55e6feab",
        "name": {
            "ar": "الليث",
            "en": "AL LITH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3428bf8d5eee6feac",
        "name": {
            "ar": "مدينة الليث",
            "en": "AL LITH TOWN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f35ce7253ee6fead",
        "name": {
            "ar": "المهد",
            "en": "AL MAHD"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f376e8138fe6feae",
        "name": {
            "ar": "المندق",
            "en": "AL MANDAQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3491727dbb4e6feaf",
        "name": {
            "ar": "المسفلة",
            "en": "AL MASFALAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38f97b8c5b9e6feb0",
        "name": {
            "ar": "الموية",
            "en": "AL MAWYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31f7e180d0de6feb1",
        "name": {
            "ar": "المخواة",
            "en": "AL MIKHWA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f33952023673e6feb2",
        "name": {
            "ar": "المويان",
            "en": "AL MOYAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3028720cfe1e6feb3",
        "name": {
            "ar": "المبرز",
            "en": "AL MUBARRAZ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f399cf32f9a2e6feb4",
        "name": {
            "ar": "المثيلف",
            "en": "AL MUTHEILIEF"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36a6e3f9c93e6feb5",
        "name": {
            "ar": "النبهانية",
            "en": "AL NABHANYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c8a8ac10bee6feb6",
        "name": {
            "ar": "النماس",
            "en": "AL NAMASS"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3240c54ce68e6feb7",
        "name": {
            "ar": "العلا",
            "en": "AL OLA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e90e74b9dde6feb8",
        "name": {
            "ar": "العتيبة",
            "en": "AL OTAIBYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c721b250dde6feb9",
        "name": {
            "ar": "العيون",
            "en": "AL OYOUN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38c62d9f4aee6feba",
        "name": {
            "ar": "القصب",
            "en": "AL QASAB"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3dfdc304716e6febb",
        "name": {
            "ar": "القصيم",
            "en": "AL QASSIM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38b8c999e6ae6febc",
        "name": {
            "ar": "الكوارع",
            "en": "AL QUAARA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f32e07d0c3fbe6febd",
        "name": {
            "ar": "الرفاعي",
            "en": "AL RAFAYE"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f64add2e2ae6febe",
        "name": {
            "ar": "الرين",
            "en": "AL RAIN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e7863cb42ae6febf",
        "name": {
            "ar": "الرس",
            "en": "AL RASS"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3a31afad4a8e6fec0",
        "name": {
            "ar": "الصحنة",
            "en": "AL SAHNA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3d678d28dade6fec1",
        "name": {
            "ar": "السيل الأكبر",
            "en": "AL SAYEL AL AKBAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f320208fd0cbe6fec2",
        "name": {
            "ar": "الشنان",
            "en": "AL SHANAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f32647cb6a6ee6fec3",
        "name": {
            "ar": "الشرفية",
            "en": "AL SHARAFIYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3211e565f63e6fec4",
        "name": {
            "ar": "الشهداء",
            "en": "AL SHUHADA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f346140cadf0e6fec5",
        "name": {
            "ar": "السحيمي",
            "en": "AL SUHEIMI"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3834483718ce6fec6",
        "name": {
            "ar": "السويتن الهوية",
            "en": "AL SUITEN AL HAWYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f35346509ddbe6fec7",
        "name": {
            "ar": "الواديين",
            "en": "AL WADYAYN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f350288d5fd5e6fec8",
        "name": {
            "ar": "الوجه",
            "en": "AL WAJH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f348f8cb95a0e6fec9",
        "name": {
            "ar": "اليتمة",
            "en": "AL YATMA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e85f91c82de6feca",
        "name": {
            "ar": "الظاهر",
            "en": "AL ZAHER"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36a48901dd3e6fecb",
        "name": {
            "ar": "أنق",
            "en": "ANAQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3a653edf06ae6fecc",
        "name": {
            "ar": "العقيق",
            "en": "AQEEQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b4acbec003e6fecd",
        "name": {
            "ar": "عرادة",
            "en": "ARADAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3659f802cd4e6fece",
        "name": {
            "ar": "عرار",
            "en": "ARAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f33b5f6ff259e6fecf",
        "name": {
            "ar": "ارطاوية",
            "en": "ARTAWIAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f39ba88887d6e6fed0",
        "name": {
            "ar": "عيون الجواء",
            "en": "AYOON AL JAWA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31093ceb53ae6fed1",
        "name": {
            "ar": "البداية",
            "en": "BADAYAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f37459ed5d76e6fed2",
        "name": {
            "ar": "بدر",
            "en": "BADR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3375331e24de6fed3",
        "name": {
            "ar": "بيش",
            "en": "BAEISH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f35d7a8fc042e6fed4",
        "name": {
            "ar": "باجاديا",
            "en": "BAGADIA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38d6f21544ae6fed5",
        "name": {
            "ar": "بها",
            "en": "BAHA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f34ade918452e6fed6",
        "name": {
            "ar": "بحرة",
            "en": "BAHRA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f37b46accfc8e6fed7",
        "name": {
            "ar": "بلجرشي",
            "en": "BALJURSHI"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f37fa0403e88e6fed8",
        "name": {
            "ar": "بالاسمار",
            "en": "BALLASMAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e8fe034514e6fed9",
        "name": {
            "ar": "بني عمار",
            "en": "BANI AMMR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3fe278ac0a8e6feda",
        "name": {
            "ar": "البقعة",
            "en": "BAQAA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f32f4b083e6be6fedb",
        "name": {
            "ar": "بارق",
            "en": "BAREQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f39eb2afaf81e6fedc",
        "name": {
            "ar": "حدود البطحاء",
            "en": "BATHA BORDER"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c2b7b1ca0fe6fedd",
        "name": {
            "ar": "بيش",
            "en": "BEISH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f371b7d4a592e6fede",
        "name": {
            "ar": "بيش",
            "en": "BISHA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f34286bf8b60e6fedf",
        "name": {
            "ar": "برازان",
            "en": "BRAZAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3cd3c41e9efe6fee0",
        "name": {
            "ar": "البكيرية",
            "en": "BUKAYRIA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e1d4bc1e7ce6fee1",
        "name": {
            "ar": "بريدة",
            "en": "BURAIDAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f32657699446e6fee2",
        "name": {
            "ar": "الدليمية",
            "en": "DALEMIYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c621c55cc8e6fee3",
        "name": {
            "ar": "داماد",
            "en": "DAMAD"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e72bc26393e6fee4",
        "name": {
            "ar": "الدمام",
            "en": "DAMMAM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f30f581d8228e6fee5",
        "name": {
            "ar": "درب",
            "en": "DARB"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3454ace42c2e6fee6",
        "name": {
            "ar": "الدوادمي",
            "en": "DAWADMI"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f33f264438a6e6fee7",
        "name": {
            "ar": "دومة الجندل",
            "en": "DAWMAT AL JANDAL"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3595490cf06e6fee8",
        "name": {
            "ar": "ذهباب",
            "en": "DHABAB"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31f94c7dbbee6fee9",
        "name": {
            "ar": "الظهران",
            "en": "DHAHRAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ade313381be6feea",
        "name": {
            "ar": "ظهران الجنوب",
            "en": "DHAHRAN AL JANOUB"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3eb26a38635e6feeb",
        "name": {
            "ar": "دوبات",
            "en": "DHOBBAT"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3282d38f449e6feec",
        "name": {
            "ar": "دوكنا",
            "en": "DHUKNA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f302fccb772fe6feed",
        "name": {
            "ar": "دولوم",
            "en": "DHULUM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31b45c9d5bae6feee",
        "name": {
            "ar": "ديلام",
            "en": "DILAM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3dfa13562b9e6feef",
        "name": {
            "ar": "دوبا",
            "en": "DUBA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f34dc13772b2e6fef0",
        "name": {
            "ar": "دورما",
            "en": "DURMA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f360675ab785e6fef1",
        "name": {
            "ar": "دوريا",
            "en": "DURYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3735021d072e6fef2",
        "name": {
            "ar": "إسكان",
            "en": "ESKAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31c77ae3731e6fef3",
        "name": {
            "ar": "فوارة",
            "en": "FAWARA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e1f8da8ea0e6fef4",
        "name": {
            "ar": "قاسم",
            "en": "GASSIM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38ab986cb65e6fef5",
        "name": {
            "ar": "جيلا",
            "en": "GELLA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f393f9538dffe6fef6",
        "name": {
            "ar": "جيزان",
            "en": "GIZAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f32f920e1351e6fef7",
        "name": {
            "ar": "القريات",
            "en": "GURAYAT"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f35f7edc82cfe6fef8",
        "name": {
            "ar": "هابونا",
            "en": "HABUNA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f390d14827e2e6fef9",
        "name": {
            "ar": "هادا",
            "en": "HADA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f344f4e5deb6e6fefa",
        "name": {
            "ar": "حفر الباطن",
            "en": "HAFAR AL BATIN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f0a34458cfe6fefb",
        "name": {
            "ar": "يشيد",
            "en": "HAIL"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3de8c3da8cde6fefc",
        "name": {
            "ar": "حلة عمار",
            "en": "HALAT AMMAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e7fd2a7ba1e6fefd",
        "name": {
            "ar": "حناكيه",
            "en": "HANAKIYAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b68ff64db2e6fefe",
        "name": {
            "ar": "حقل",
            "en": "HAQL"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3330b7e4c0de6feff",
        "name": {
            "ar": "حرض",
            "en": "HARADH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f39599bd7b42e6ff00",
        "name": {
            "ar": "حودات سدير",
            "en": "HAWDAT SUDAIR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f00ff5a49fe6ff01",
        "name": {
            "ar": "حودات تميم",
            "en": "HAWDAT TAMIM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3a882f79337e6ff02",
        "name": {
            "ar": "حياة",
            "en": "HAYET"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38d751c762be6ff03",
        "name": {
            "ar": "حائر",
            "en": "HAYIR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f383b132bf3be6ff04",
        "name": {
            "ar": "حلبان",
            "en": "HELBAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36ee7b6353ee6ff05",
        "name": {
            "ar": "الحجاز",
            "en": "HIJAZ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f39003054a3de6ff06",
        "name": {
            "ar": "هجرة لبن",
            "en": "HIJRAT LABAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f086d676ede6ff07",
        "name": {
            "ar": "هيتم",
            "en": "HITEEM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b6aaf07547e6ff08",
        "name": {
            "ar": "الهفوف",
            "en": "HOFUF"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e88063ac3fe6ff09",
        "name": {
            "ar": "حطة بن تميم",
            "en": "HOTA BIN TAMIM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f358073fdcb3e6ff0a",
        "name": {
            "ar": "هوتا سودهير",
            "en": "HOTA SUDHAIR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3aa88e2918fe6ff0b",
        "name": {
            "ar": "حريمالا",
            "en": "HURAIMALA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f34523d545aee6ff0c",
        "name": {
            "ar": "جامعة الامام",
            "en": "IMAM UNIVERSITY"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f34e6951f014e6ff0d",
        "name": {
            "ar": "جبل النور",
            "en": "JABAL AL NOOR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e4a8d2dc4ae6ff0e",
        "name": {
            "ar": "جعفر",
            "en": "JAFAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36b854e01ffe6ff0f",
        "name": {
            "ar": "جلاجل",
            "en": "JALAJIL"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f316179f554ae6ff10",
        "name": {
            "ar": "جمجوم",
            "en": "JAMJOOM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3d5693806aee6ff11",
        "name": {
            "ar": "جاردا",
            "en": "JARDA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f39474c2f954e6ff12",
        "name": {
            "ar": "جرير",
            "en": "JAREER"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c1b04860bfe6ff13",
        "name": {
            "ar": "جدة",
            "en": "JEDDAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f364319fdf3be6ff14",
        "name": {
            "ar": "الجبيل",
            "en": "JUBAIL"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f39ea9a63f82e6ff15",
        "name": {
            "ar": "جبة",
            "en": "JUBBAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f35f56aff4f5e6ff16",
        "name": {
            "ar": "الجهيمية",
            "en": "JUHAIMIA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f30032e0cd6ae6ff17",
        "name": {
            "ar": "كفى",
            "en": "KAFA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c3b30b39cae6ff18",
        "name": {
            "ar": "خبرا",
            "en": "KHABRA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f367bae6498fe6ff19",
        "name": {
            "ar": "الخفجي",
            "en": "KHAFJI"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3172b1e3a22e6ff1a",
        "name": {
            "ar": "خيبر",
            "en": "KHAIBAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c338cd6d3ce6ff1b",
        "name": {
            "ar": "الخالدية",
            "en": "KHALIDYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3100de03a82e6ff1c",
        "name": {
            "ar": "خميس مشيط",
            "en": "KHAMIS MUSHAYAT"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f30c3ef58dd1e6ff1d",
        "name": {
            "ar": "جيزان",
            "en": "KHAZZAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f300d9f5b24de6ff1e",
        "name": {
            "ar": "خبيب",
            "en": "KHUBAIB"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b9a9c23ff1e6ff1f",
        "name": {
            "ar": "خريص",
            "en": "KHURAIS"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f33f574e13d6e6ff20",
        "name": {
            "ar": "خورما",
            "en": "KHURMA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ab4e8c0eefe6ff21",
        "name": {
            "ar": "خطا",
            "en": "KHUTTA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38182175851e6ff22",
        "name": {
            "ar": "محايل عسير",
            "en": "MAHAYIL ASEER"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38a5c6e090ae6ff23",
        "name": {
            "ar": "جاردة",
            "en": "MAJARDEH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3d779b11db7e6ff24",
        "name": {
            "ar": "مجمعة",
            "en": "MAJMAAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f37e1d48dcc5e6ff25",
        "name": {
            "ar": "مكه",
            "en": "MAKKAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f34a814b6f47e6ff26",
        "name": {
            "ar": "مندق",
            "en": "MANDAQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e7bba57a9ce6ff27",
        "name": {
            "ar": "منفوحة",
            "en": "MANFOUAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3510fa4a2c6e6ff28",
        "name": {
            "ar": "منفوحة",
            "en": "MANFOUHA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e57d098f80e6ff29",
        "name": {
            "ar": "مانيفا",
            "en": "MANIFA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3211af62a90e6ff2a",
        "name": {
            "ar": "ام",
            "en": "MATHER"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ca6e2b57b1e6ff2b",
        "name": {
            "ar": "موقع",
            "en": "MAWQEQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3fefb4b941fe6ff2c",
        "name": {
            "ar": "المزروعية",
            "en": "MAZRUEYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f380f2b43367e6ff2d",
        "name": {
            "ar": "المدينة المنورة",
            "en": "MEDINA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f39ae4be13f8e6ff2e",
        "name": {
            "ar": "ميجايبراه",
            "en": "MEGAIBRAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f6459bb6a4e6ff2f",
        "name": {
            "ar": "مدنب",
            "en": "MIDHNAB"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31558361e29e6ff30",
        "name": {
            "ar": "مغرزات",
            "en": "MOGHARAZAT"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38e473f6aa4e6ff31",
        "name": {
            "ar": "موكاك",
            "en": "MOQAQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3db1aa8e9a7e6ff32",
        "name": {
            "ar": "المبرز",
            "en": "MUBARRAZ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f361eca15fd3e6ff33",
        "name": {
            "ar": "مناخة",
            "en": "MUNAKHAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e0f18de767e6ff34",
        "name": {
            "ar": "مؤتمرات",
            "en": "MUTAMARAT"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c94c3ba631e6ff35",
        "name": {
            "ar": "مزامية",
            "en": "MUZAMMIA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f381d2da3f8ae6ff36",
        "name": {
            "ar": "النبانية",
            "en": "NABANIYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c5104025cbe6ff37",
        "name": {
            "ar": "نافا",
            "en": "NAFA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f367549ffc2ae6ff38",
        "name": {
            "ar": "نعيم",
            "en": "NAIM"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31cb38a534fe6ff39",
        "name": {
            "ar": "ناريا",
            "en": "NARIYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3681dc61bd9e6ff3a",
        "name": {
            "ar": "نجران",
            "en": "NEJRAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b98335502ae6ff3b",
        "name": {
            "ar": "الناصرية",
            "en": "NESRIYAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f30ae17f0faee6ff3c",
        "name": {
            "ar": "الشميسي الجديد",
            "en": "NEW SHIMAISY"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f35e8477cf2ae6ff3d",
        "name": {
            "ar": "عقلة الصقور",
            "en": "OQLAH AL SUQOUR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f33231cb0751e6ff3e",
        "name": {
            "ar": "أوشايجر",
            "en": "OSHAIGER"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f37aec7e637ee6ff3f",
        "name": {
            "ar": "عيون",
            "en": "OYOON"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f362a03d062de6ff40",
        "name": {
            "ar": "القيصومة",
            "en": "QAISUMAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f32625f17b4be6ff41",
        "name": {
            "ar": "قرية العليا",
            "en": "QARYAT AL OLAYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3a4b64a402ce6ff42",
        "name": {
            "ar": "القطيف",
            "en": "QATIF"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3060561a727e6ff43",
        "name": {
            "ar": "قلوة",
            "en": "QILWAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36c51c9fce9e6ff44",
        "name": {
            "ar": "قوبا",
            "en": "QUBA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3d9e6ea03c3e6ff45",
        "name": {
            "ar": "قبة",
            "en": "QUBBA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3002591aac8e6ff46",
        "name": {
            "ar": "القنفذة",
            "en": "QUNFUDAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c0e3825a65e6ff47",
        "name": {
            "ar": "القويعية",
            "en": "QUWAYIAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3fc8247ca88e6ff48",
        "name": {
            "ar": "رابغ",
            "en": "RABIGH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3adcee19bdde6ff49",
        "name": {
            "ar": "رفحاء",
            "en": "RAFHA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f32aa770c3e4e6ff4a",
        "name": {
            "ar": "رحيمة",
            "en": "RAHIMAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e3b25da2b1e6ff4b",
        "name": {
            "ar": "رهوة البر",
            "en": "RAHWA AL BAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3095f32f024e6ff4c",
        "name": {
            "ar": "رانيا",
            "en": "RANIA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f378f5f2e077e6ff4d",
        "name": {
            "ar": "رأس تنورة",
            "en": "RAS TANURA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36b21a9f731e6ff4e",
        "name": {
            "ar": "رجال ألما",
            "en": "REJAL ALMAA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3a8104e4a57e6ff4f",
        "name": {
            "ar": "الرياض",
            "en": "RIYADH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ba280cb1bee6ff50",
        "name": {
            "ar": "الرياض الخبراء",
            "en": "RIYADH ALKHBRA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b7069b47d6e6ff51",
        "name": {
            "ar": "روضة",
            "en": "ROWADA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31b1ab2b08fe6ff52",
        "name": {
            "ar": "الرويضة",
            "en": "ROWAIDA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3281428bf41e6ff53",
        "name": {
            "ar": "روماه",
            "en": "RUMAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3397569a687e6ff54",
        "name": {
            "ar": "الرويضة",
            "en": "RUWAIDAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f339489c935de6ff55",
        "name": {
            "ar": "سبتاليا",
            "en": "SABTALAIA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f328206c1acee6ff56",
        "name": {
            "ar": "صبيا",
            "en": "SABYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c326c614e2e6ff57",
        "name": {
            "ar": "صفاء",
            "en": "SAFA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3eeb3b81419e6ff58",
        "name": {
            "ar": "السفانية",
            "en": "SAFANIYAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ede7090ad3e6ff59",
        "name": {
            "ar": "صفرا",
            "en": "SAFRA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f344555a1170e6ff5a",
        "name": {
            "ar": "الصفوة",
            "en": "SAFWA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3d292f9af8ae6ff5b",
        "name": {
            "ar": "سيهات",
            "en": "SAIHAT"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f349b80c7a63e6ff5c",
        "name": {
            "ar": "ساجير",
            "en": "SAJIR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f78bd2540de6ff5d",
        "name": {
            "ar": "سكاكا",
            "en": "SAKAKAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ba1cb0c915e6ff5e",
        "name": {
            "ar": "سلامة",
            "en": "SALAMAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f100352520e6ff5f",
        "name": {
            "ar": "سلوى",
            "en": "SALWA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f387612ce841e6ff60",
        "name": {
            "ar": "سماح",
            "en": "SAMAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f39b6a7a515fe6ff61",
        "name": {
            "ar": "سماشيا",
            "en": "SAMASHIYA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36dc695848de6ff62",
        "name": {
            "ar": "صامطة",
            "en": "SAMTAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36c16b63a5ae6ff63",
        "name": {
            "ar": "سارار",
            "en": "SARAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38343f77717e6ff64",
        "name": {
            "ar": "سراة عبيدة",
            "en": "SARAT ABIDAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3351581daf9e6ff65",
        "name": {
            "ar": "شقرا",
            "en": "SHAHQRA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f341e850f9a2e6ff66",
        "name": {
            "ar": "شمالي",
            "en": "SHAMALI"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3e4fd8ba1c7e6ff67",
        "name": {
            "ar": "شمسان",
            "en": "SHAMSAN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f373de7c530be6ff68",
        "name": {
            "ar": "شاري",
            "en": "SHARI"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3746efd8501e6ff69",
        "name": {
            "ar": "شرورة",
            "en": "SHARURAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f349bd7f2892e6ff6a",
        "name": {
            "ar": "الشعيبة",
            "en": "SHOAIBA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f32e81d3ae74e6ff6b",
        "name": {
            "ar": "سياهات",
            "en": "SIAHAT"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38e09960422e6ff6c",
        "name": {
            "ar": "سيتين",
            "en": "SITTEN"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3d8cf05596be6ff6d",
        "name": {
            "ar": "سوق الأحد",
            "en": "SOUQ ALAHAD"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f31b43f3d83be6ff6e",
        "name": {
            "ar": "السليل",
            "en": "SULAYYIL"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3197109c72be6ff6f",
        "name": {
            "ar": "سميراء",
            "en": "SUMEIRA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3168576732ce6ff70",
        "name": {
            "ar": "تبوك",
            "en": "TABUK"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3984e8fc0f8e6ff71",
        "name": {
            "ar": "طائف",
            "en": "TAIF"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f377f02e3d79e6ff72",
        "name": {
            "ar": "تانداها",
            "en": "TANDAHA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3aa30768464e6ff73",
        "name": {
            "ar": "تنوما",
            "en": "TANUMA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3cbf37b3328e6ff74",
        "name": {
            "ar": "جزيرة تاروت",
            "en": "TARUT ISLAND"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3dbc8404c1de6ff75",
        "name": {
            "ar": "تعاون",
            "en": "TAWOON"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f33967de22e2e6ff76",
        "name": {
            "ar": "تيمة",
            "en": "TAYMAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f33bbd90c28de6ff77",
        "name": {
            "ar": "ثادق",
            "en": "THADIQ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38b5f073cc8e6ff78",
        "name": {
            "ar": "ثريب",
            "en": "THAREEB"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f38fa0b82802e6ff79",
        "name": {
            "ar": "ثاتليت",
            "en": "THATLEETH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3ad4737cd2ce6ff7a",
        "name": {
            "ar": "ثومير",
            "en": "THOMAIR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f392e579453de6ff7b",
        "name": {
            "ar": "ثوال",
            "en": "THOWAL"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36f05b0293be6ff7c",
        "name": {
            "ar": "توجاباه",
            "en": "TOUGABAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f352d3ee9d04e6ff7d",
        "name": {
            "ar": "طريف",
            "en": "TURAIF"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f348ff12bcf9e6ff7e",
        "name": {
            "ar": "تربة الشمال",
            "en": "TURBAH AL SHAMAL"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f360cdf64f5fe6ff7f",
        "name": {
            "ar": "العديلية",
            "en": "UDILLIYAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f33d376b221be6ff80",
        "name": {
            "ar": "ام القرى",
            "en": "UM AL QURA"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f36a14dee51fe6ff81",
        "name": {
            "ar": "ام الصادق",
            "en": "UM AL SAHEK"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3f1b8e2be8ee6ff82",
        "name": {
            "ar": "أملج",
            "en": "UMLUJJ"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3c1dfbfcfd2e6ff83",
        "name": {
            "ar": "عنيزة",
            "en": "UNAYZAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f397f256b5d8e6ff84",
        "name": {
            "ar": "عقلاق الصقر",
            "en": "UQLAQ AL SUGGUR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f32d4f68e056e6ff85",
        "name": {
            "ar": "يوينا",
            "en": "UYENAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f35acc0ad09ee6ff86",
        "name": {
            "ar": "وادي الدواسر",
            "en": "WADI AL DAWASIR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f358ae67abd6e6ff87",
        "name": {
            "ar": "وادي هشبل",
            "en": "WADI HASHBEL"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3a00a715a02e6ff88",
        "name": {
            "ar": "ورود",
            "en": "WOROOD"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b995c2d364e6ff89",
        "name": {
            "ar": "ينبوع",
            "en": "YANBU"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f358249b71ace6ff8a",
        "name": {
            "ar": "الزاوية",
            "en": "ZAWIYAH"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f30861c0a0eee6ff8b",
        "name": {
            "ar": "زديهار",
            "en": "ZDIHAR"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f3b233f77a8fe6ff8c",
        "name": {
            "ar": "زلفي",
            "en": "ZILFI"
        },
		"country_id" : "SA"
    },
    {
        "_id": "6666c1f30861c0a0eee6ff6e",
        "name": {
            "ar": "القاهرة",
            "en": "Cairo"
        },
		"country_id" : "EG"
    }
]`

const Countries = `[
    {
        "_id": "SA",
        "name": {
            "ar": "المملكة العربية السعودية",
            "en": "Saudi Arabia"
        },
	"timezone" : "Asia/Riyadh",
	"currency" : "SAR",
	"phone_prefix" : "966"
    },
    {
        "_id": "EG",
        "name": {
            "ar": "جمهورية مصر العربية",
            "en": "Egypt"
        },
	"timezone" : "Asia/Riyadh",
	"currency" : "EGP",
	"phone_prefix" : "20"
    }
]`

type Name struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

// City represents a city with an ID and name
type City struct {
	ID        string `json:"_id"`
	Name      Name   `json:"name"`
	CountryId string `json:"country_id"`
}

func CitiesBuilder(payload *dto.ListCitiesDto) interface{} {
	var data []City

	err := json.Unmarshal([]byte(Cities), &data)
	if err != nil {
		return data
	}
	result := data
	if payload.CountryId != "" {
		From(data).Where(func(c interface{}) bool {
			return c.(City).CountryId == strings.ToUpper(payload.CountryId)
		}).ToSlice(&result)
	}

	return result
}
func CountriesBuilder() interface{} {
	var data interface{}

	err := json.Unmarshal([]byte(Countries), &data)
	if err != nil {
		fmt.Println("error => ", err)
		return data
	}
	return data
}

func ReadFile(iLogger logger.ILogger, filePath string) []map[string]interface{} {
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + filePath)
	if err != nil {
		iLogger.Error(err)
	}

	defer file.Close()
	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		iLogger.Error(err)
	}

	var result []map[string]interface{}
	json.Unmarshal(content, &result)
	return result
}
