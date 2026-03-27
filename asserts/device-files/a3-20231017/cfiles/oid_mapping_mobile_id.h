#ifndef __OID_MAPPING_MOBILE_ID_H__
#define __OID_MAPPING_MOBILE_ID_H__
#include "sunwave_oid.h"
#include "../Snmp/SunwaveMib.h"

#define ALM_THR_ROOT    (2)
#define LAN_SET_ROOT    (4)
#define SYS_INFO_ROOT   (5)
#define BAND_CONF_ROOT  (7)
#define POI_ROOT        (8)
#define GAIN_ROOT       (9)
#define CAPACITY_ROOT   (10)
/*last 3 bytes of scalar variables oid*/
/**************告警门限*******************/
#define TEMP_HIGH_THR_OID           {2,2,1}

/*Au*/
#define DN_INPWR_LOW_THR1_OID       {2,2,2}
#define DN_INPWR_LOW_THR2_OID       {2,2,3}
#define DN_INPWR_LOW_THR3_OID       {2,2,4}
#define DN_INPWR_LOW_THR4_OID       {2,2,5}

#define DN_INPWR_HIGH_THR1_OID      {2,2,6}
#define DN_INPWR_HIGH_THR2_OID      {2,2,7}
#define DN_INPWR_HIGH_THR3_OID      {2,2,8}
#define DN_INPWR_HIGH_THR4_OID      {2,2,9}

#define EXT_ALM1_MODE_SEL      		{2,3,1}
#define EXT_ALM2_MODE_SEL      		{2,3,2}
#define EXT_ALM3_MODE_SEL      		{2,3,3}
#define EXT_ALM4_MODE_SEL      		{2,3,4}

#define EXT_ALM1_REMARK      		{2,3,5}
#define EXT_ALM2_REMARK      		{2,3,6}
#define EXT_ALM3_REMARK      		{2,3,7}
#define EXT_ALM4_REMARK      		{2,3,8}

#define EXT_ALM_OUT_WARNING_MODE_SEL	{2,3,9}
#define EXT_ALM_OUT_MINOR_MODE_SEL     	{2,3,10}
#define EXT_ALM_OUT_MAJOR_MODE_SEL		{2,3,11}
#define EXT_ALM_OUT_CRITICAL_MODE_SEL	{2,3,12}

#define EXT_ALM_OUT_WARNING_STATE      	{2,3,13}
#define EXT_ALM_OUT_MINOR_STATE     	{2,3,14}
#define EXT_ALM_OUT_MAJOR_STATE      	{2,3,15}
#define EXT_ALM_OUT_CRITICAL_STATE      {2,3,16}

#define EXT_ALM_OUT_1_MODE_SEL 			{2,3,17}
#define EXT_ALM_OUT_2_MODE_SEL			{2,3,18}
#define EXT_ALM_OUT_1_Level				{2,3,19}
#define EXT_ALM_OUT_2_Level				{2,3,20}
#define EXT_ALM_OUT_1_STATE 			{2,3,21}
#define EXT_ALM_OUT_2_STATE				{2,3,22}

#define EXT_OUT_ALM1_MODE_SEL      		{2,3,23}
#define EXT_OUT_ALM2_MODE_SEL      		{2,3,24}

#define EXT_OUT_ALM1_LEVEL_SET      	{2,3,25}
#define EXT_OUT_ALM2_LEVEL_SET      	{2,3,26}

#define EXT_ALM1_OUT_INDICATION      	{2,3,27}
#define EXT_ALM2_OUT_INDICATION      	{2,3,28}




/*RU*/
#define DN_OUTPWR_LOW_THR1_OID      {2,2,2}
#define DN_OUTPWR_LOW_THR2_OID      {2,2,3}
#define DN_OUTPWR_LOW_THR3_OID      {2,2,4}
#define DN_OUTPWR_LOW_THR4_OID      {2,2,5}
/*N3RU*/
#define DN_OUTPWR_LOW_THR5_OID      {2,2,27}
#define DN_OUTPWR_LOW_THR6_OID      {2,2,28}
#define DN_OUTPWR_LOW_THR7_OID      {2,2,29}
#define DN_OUTPWR_LOW_THR8_OID      {2,2,30}

#define DN_OUTPWR_HIGH_THR1_OID     {2,2,6}
#define DN_OUTPWR_HIGH_THR2_OID     {2,2,7}
#define DN_OUTPWR_HIGH_THR3_OID     {2,2,8}
#define DN_OUTPWR_HIGH_THR4_OID     {2,2,9}
/*N3RU*/
#define DN_OUTPWR_HIGH_THR5_OID     {2,2,31}
#define DN_OUTPWR_HIGH_THR6_OID     {2,2,32}
#define DN_OUTPWR_HIGH_THR7_OID     {2,2,33}
#define DN_OUTPWR_HIGH_THR8_OID     {2,2,34}


/*HRU*/
#define VSWR_THR_OID                {2,2,10}

/*N2RU*/
#define N2RU_DN_OUTPWR_LOW_THR1_OID      {2,2,11}
#define N2RU_DN_OUTPWR_LOW_THR2_OID      {2,2,12}
#define N2RU_DN_OUTPWR_LOW_THR3_OID      {2,2,13}
#define N2RU_DN_OUTPWR_LOW_THR4_OID      {2,2,14}
#define N2RU_DN_OUTPWR_LOW_THR5_OID      {2,2,15}
#define N2RU_DN_OUTPWR_LOW_THR6_OID      {2,2,16}
#define N2RU_DN_OUTPWR_LOW_THR7_OID      {2,2,17}
#define N2RU_DN_OUTPWR_LOW_THR8_OID      {2,2,18}

#define N2RU_DN_OUTPWR_HIGH_THR1_OID     {2,2,19}
#define N2RU_DN_OUTPWR_HIGH_THR2_OID     {2,2,20}
#define N2RU_DN_OUTPWR_HIGH_THR3_OID     {2,2,21}
#define N2RU_DN_OUTPWR_HIGH_THR4_OID     {2,2,22}
#define N2RU_DN_OUTPWR_HIGH_THR5_OID     {2,2,23}
#define N2RU_DN_OUTPWR_HIGH_THR6_OID     {2,2,24}
#define N2RU_DN_OUTPWR_HIGH_THR7_OID     {2,2,25}
#define N2RU_DN_OUTPWR_HIGH_THR8_OID     {2,2,26}
/*******************************************/

/******************设备基本信息*************/
#define FACTORY_ID_OID              {4,1,1}
#define DEV_MODEL_OID               {4,1,2}
#define DEV_SRRIAL_OID              {4,1,3}
#define ARM_VER_OID                 {4,1,4}
#define SNMP_VER_OID                {4,1,5}
#define DEVICE_ID                   {4,1,6}
#define DEVICE_SUB_ID               {4,1,7}
#define DEV_LOCATION                {4,1,8}
#define DEVICE_NAME                 {4,1,9}
#define DEV_TIME_OID                {4,1,10}
#define SITE_REPORT_OID             {4,1,14}
#define SITE_NAME             		{4,1,15}
#define MAC_ADDRESS             	{4,1,30}
#define LIFE_TIME	             	{4,1,31}
#define UP_TIME	             		{4,1,32}
#define SYSTEM_RESET	            {4,1,33}
#define FILESYSTEM_VER_OID	        {4,1,34}


/*MasterAu*/
#define PROTOCOL_SEL                {4,2,1}
#define NMS_IPADDR                  {4,2,2}
#define NMS_PORT                    {4,2,3} 
#define DEVICE_IPADD_OID            {4,2,7}
#define NET_MASK_ADD_OID            {4,2,8}
#define DEFAULT_GWAY_OID            {4,2,9}
#define DEVREV_PORT                 {4,2,10}
#define HEARTTIME                   {4,2,11} 
#define DEVICE_IPADD2_OID           {4,2,12}
#define NET_MASK2_ADD_OID           {4,2,13}
#define DEFAULT_GWAY2_OID           {4,2,14}


#define FTP_IP_OID                  {4,3,1}
#define FTP_PORT_OID                {4,3,2}
#define FTP_USER_OID                {4,3,3}
#define FTP_PASSRD_OID              {4,3,4}
#define FTP_DIR_OID                 {4,3,5}
#define FTP_FILENAME_OID            {4,3,6}
#define FTP_TRANS_OID               {4,3,7}

#define SEC_USER_NAME_OID           {4,4,1}
#define AUTH_PROTOCOL_OID           {4,4,2}
#define AUTH_PASSWORD_OID           {4,4,3}
#define PRIV_PROTOCOL_OID           {4,4,4}
#define PRIV_PASSWORD_OID           {4,4,5}
#define USM_EDIT_COMFIRM_OID        {4,4,6}
#define PROTOCOL_SEL_RO_OID         {4,4,7}
#define SNMP_TRAP_PROTOCOL          {4,4,8}
#define USM_USR_UPDATE_OID          {4,4,9}
#define USM_CURRENT_USR_OID         {4,4,10}
#define USM_DELETE_USR_OID          {4,4,11}


//#define TRAP_IPADD_OID              {4,4,9}
//#define TRAP_IPADD2_OID             {4,4,10}
//#define TRAP_PORT_OID               {4,4,11}
//#define TRAP_ENGINE1_ID             {4,4,12}
//#define TRAP_ENGINE2_ID             {4,4,13}
#define USM_RESET_OID               {4,4,14}

#define RESEND_ENABLE_OID      		{4,4,15}
#define RESEND_INTERVAL_OID   		{4,4,16}
#define COMMUNITY_OID         	    {4,4,17}
#define DELETE_HIS_AlM_OID      	{4,4,18}

#define TRAP_IPADD_OID              {4,4,19}
#define TRAP_IPADD2_OID             {4,4,20}
#define TRAP_IPADD3_OID             {4,4,21}
#define TRAP_IPADD4_OID             {4,4,22}
#define TRAP_IPADD5_OID             {4,4,23}
#define TRAP_IPADD6_OID             {4,4,24}
#define TRAP_IPADD7_OID             {4,4,25}
#define TRAP_IPADD8_OID             {4,4,26}
#define TRAP_IPADD9_OID             {4,4,27}
#define TRAP_IPADD10_OID            {4,4,28}
#define TRAP_PORT_OID               {4,4,29}

#define TRAP_ENGINE1_ID             {4,4,30}
#define TRAP_ENGINE2_ID             {4,4,31}
#define TRAP_ENGINE3_ID             {4,4,32}
#define TRAP_ENGINE4_ID             {4,4,33}
#define TRAP_ENGINE5_ID             {4,4,34}
#define TRAP_ENGINE6_ID             {4,4,35}
#define TRAP_ENGINE7_ID             {4,4,36}
#define TRAP_ENGINE8_ID             {4,4,37}
#define TRAP_ENGINE9_ID             {4,4,38}
#define TRAP_ENGINE10_ID            {4,4,39}

#define BEACON_SWITCH			    {4,5,1}
#define BEACON_UUID			        {4,5,2}
#define BEACON_MAJOR			    {4,5,3}
#define BEACON_MINOR			    {4,5,4}

#define SEC_USER_NAME1_OID         {4,5,15}
#define AUTH_PROTOCOL1_OID         {4,5,16}
#define AUTH_PASSWORD1_OID         {4,5,17}
#define PRIV_PROTOCOL1_OID         {4,5,18}
#define PRIV_PASSWORD1_OID         {4,5,19}
#define USM_EDIT_COMFIRM1_OID      {4,5,20}
//#define TRAP_IP1ADD_OID            {4,5,20}

#define SEC_USER_NAME2_OID         {4,5,21}
#define AUTH_PROTOCOL2_OID         {4,5,22}
#define AUTH_PASSWORD2_OID         {4,5,23}
#define PRIV_PROTOCOL2_OID         {4,5,24}
#define PRIV_PASSWORD2_OID         {4,5,25}
#define USM_EDIT_COMFIRM2_OID      {4,5,26}
//#define TRAP_IP2ADD_OID            {4,5,26} 

#define SEC_USER_NAME3_OID         {4,5,27}
#define AUTH_PROTOCOL3_OID         {4,5,28}
#define AUTH_PASSWORD3_OID         {4,5,29}
#define PRIV_PROTOCOL3_OID         {4,5,30}
#define PRIV_PASSWORD3_OID         {4,5,31}
#define USM_EDIT_COMFIRM3_OID      {4,5,32}
//#define TRAP_IP3ADD_OID            {4,5,32} 

#define SEC_USER_NAME4_OID         {4,5,33}
#define AUTH_PROTOCOL4_OID         {4,5,34}
#define AUTH_PASSWORD4_OID         {4,5,35}
#define PRIV_PROTOCOL4_OID         {4,5,36}
#define PRIV_PASSWORD4_OID         {4,5,37}
#define USM_EDIT_COMFIRM4_OID      {4,5,38}
//#define TRAP_IP4ADD_OID            {4,5,38} 

#define SEC_USER_NAME5_OID         {4,5,39}
#define AUTH_PROTOCOL5_OID         {4,5,40}
#define AUTH_PASSWORD5_OID         {4,5,41}
#define PRIV_PROTOCOL5_OID         {4,5,42}
#define PRIV_PASSWORD5_OID         {4,5,43}
#define USM_EDIT_COMFIRM5_OID      {4,5,44}
//#define TRAP_IP5ADD_OID            {4,5,44} 

#define SEC_USER_NAME6_OID         {4,5,45}
#define AUTH_PROTOCOL6_OID         {4,5,46}
#define AUTH_PASSWORD6_OID         {4,5,47}
#define PRIV_PROTOCOL6_OID         {4,5,48}
#define PRIV_PASSWORD6_OID         {4,5,49}
#define USM_EDIT_COMFIRM6_OID      {4,5,50}
//#define TRAP_IP6ADD_OID		       {4,5,50} 

#define SEC_USER_NAME7_OID         {4,5,51}
#define AUTH_PROTOCOL7_OID         {4,5,52}
#define AUTH_PASSWORD7_OID         {4,5,53}
#define PRIV_PROTOCOL7_OID         {4,5,54}
#define PRIV_PASSWORD7_OID         {4,5,55}
#define USM_EDIT_COMFIRM7_OID      {4,5,56}
//#define TRAP_IP7ADD_OID            {4,5,56} 

#define SEC_USER_NAME8_OID         {4,5,57}
#define AUTH_PROTOCOL8_OID         {4,5,58}
#define AUTH_PASSWORD8_OID         {4,5,59}
#define PRIV_PROTOCOL8_OID         {4,5,60}
#define PRIV_PASSWORD8_OID         {4,5,61}
#define USM_EDIT_COMFIRM8_OID      {4,5,62}
//#define TRAP_IP8ADD_OID            {4,5,62} 

#define SEC_USER_NAME9_OID         {4,5,63}
#define AUTH_PROTOCOL9_OID         {4,5,64}
#define AUTH_PASSWORD9_OID         {4,5,65}
#define PRIV_PROTOCOL9_OID         {4,5,66}
#define PRIV_PASSWORD9_OID         {4,5,67}
#define USM_EDIT_COMFIRM9_OID      {4,5,68}
//#define TRAP_IP9ADD_OID            {4,5,68} 

#define SEC_USER_NAME10_OID         {4,5,69}
#define AUTH_PROTOCOL10_OID         {4,5,70}
#define AUTH_PASSWORD10_OID         {4,5,71}
#define PRIV_PROTOCOL10_OID         {4,5,72}
#define PRIV_PASSWORD10_OID         {4,5,73}
#define USM_EDIT_COMFIRM10_OID      {4,5,74}
//#define TRAP_IP10ADD_OID            {4,5,74} 

#define SLV_DEV_NUM_OID             {1,1}
/***********************************************/

/*******************系统信息********************/
/*Au&Ru*/
#define RF_SWITCH1_OID              {5,1,1}
#define UP_GAIN_OFFSET1_OID         {5,1,2}
#define DN_GAIN_OFFSET1_OID         {5,1,3}
#define BAND_WIDE1_OID              {5,1,4} 
#define UP_CENTER_FRE1_OID          {5,1,5}
#define DN_CENTER_FRE1_OID          {5,1,6}
/*Ru*/
#define UP_INPUT1_POWER_OID         {5,1,7}
#define PA_OUTPWR1_OID              {5,1,8}
#define CH1_VSWR_OID                {5,1,9}

#define RF_SWITCH2_OID              {5,2,1}
#define UP_GAIN_OFFSET2_OID         {5,2,2}
#define DN_GAIN_OFFSET2_OID         {5,2,3}
#define BAND_WIDE2_OID              {5,2,4} 
#define UP_CENTER_FRE2_OID          {5,2,5}
#define DN_CENTER_FRE2_OID          {5,2,6}
/*Ru*/
#define UP_INPUT2_POWER_OID         {5,2,7}
#define PA_OUTPWR2_OID              {5,2,8}
#define CH2_VSWR_OID                {5,2,9}

#define RF_SWITCH3_OID              {5,3,1}
#define UP_GAIN_OFFSET3_OID         {5,3,2}
#define DN_GAIN_OFFSET3_OID         {5,3,3}
#define BAND_WIDE3_OID              {5,3,4} 
#define UP_CENTER_FRE3_OID          {5,3,5}
#define DN_CENTER_FRE3_OID          {5,3,6}
/*Ru*/
#define UP_INPUT3_POWER_OID         {5,3,7}
#define PA_OUTPWR3_OID              {5,3,8}
#define CH3_VSWR_OID                {5,3,9}

#define RF_SWITCH4_OID              {5,4,1}
#define UP_GAIN_OFFSET4_OID         {5,4,2}
#define DN_GAIN_OFFSET4_OID         {5,4,3}
#define BAND_WIDE4_OID              {5,4,4}
#define UP_CENTER_FRE4_OID          {5,4,5}
#define DN_CENTER_FRE4_OID          {5,4,6}
/*Ru*/
#define UP_INPUT4_POWER_OID         {5,4,7}
#define PA_OUTPWR4_OID              {5,4,8}
#define CH4_VSWR_OID                {5,4,9}


#define DEV_TEMP_MAX_OID            {5,5,1}
#define ROUTE_ADD_OID               {5,5,2}
#define Open_Load_Detect_Switch  	{5,5,3}

/*m2ru*/

#define CH1_BASEBAND_UL_INPUT_POWER_OID  {5,1,10}
#define CH2_BASEBAND_UL_INPUT_POWER_OID  {5,2,10}
#define CH3_BASEBAND_UL_INPUT_POWER_OID  {5,3,10}
#define CH4_BASEBAND_UL_INPUT_POWER_OID  {5,4,10}

#define CH1_BASEBAND_DL_OUTPUT_POWER_OID {5,1,11}
#define CH2_BASEBAND_DL_OUTPUT_POWER_OID {5,2,11}
#define CH3_BASEBAND_DL_OUTPUT_POWER_OID {5,3,11}
#define CH4_BASEBAND_DL_OUTPUT_POWER_OID {5,4,11}


/*SU*/
#define CH1_INPUT_POWER_OID	{5, 6, 1}
#define CH2_INPUT_POWER_OID	{5, 6, 2}
#define CH3_INPUT_POWER_OID	{5, 6, 3}
#define CH4_INPUT_POWER_OID	{5, 6, 4}

#define SHIELD_SIG_CYC_OID		{5, 7, 1}
#define LTE_SYNC_FLAG_OID		{5, 7, 2}

#define GPS_MODE_OID			{5, 8, 1}
#define AD_MODE_SELECT_OID		{5, 8, 2}

#define SET_DELAY_TIME_OID		{5, 9, 1}
#define GPS_DELAY_TIME_OID		{5, 9, 2}
#define AU_DELAY_TIME_OID		{5, 9, 3}

/*XP-RU*/
/*#define ALM_FAN_DEV_40W		{5,10,1}
#define FAN1_SPEED_40W			{5,10,2}
#define FAN2_SPEED_40W			{5,10,3}
#define FAN_HIGH_THRESHOLD	{5,10,4}
#define FAN_LOW_THRESHOLD		{5,10,5}*/

/* N2RU */
#define RF_SWITCH5_OID              {5,11,1}
#define UP_GAIN_OFFSET5_OID         {5,11,2}
#define DN_GAIN_OFFSET5_OID         {5,11,3}
#define BAND_WIDE5_OID              {5,11,4}
#define UP_CENTER_FRE5_OID          {5,11,5}
#define DN_CENTER_FRE5_OID          {5,11,6}
#define UP_INPUT5_POWER_OID         {5,11,7}
#define PA_OUTPWR5_OID              {5,11,8}
#define CH5_VSWR_OID                {5,11,9}

#define RF_SWITCH6_OID              {5,12,1}
#define UP_GAIN_OFFSET6_OID         {5,12,2}
#define DN_GAIN_OFFSET6_OID         {5,12,3}
#define BAND_WIDE6_OID              {5,12,4}
#define UP_CENTER_FRE6_OID          {5,12,5}
#define DN_CENTER_FRE6_OID          {5,12,6}
#define UP_INPUT6_POWER_OID         {5,12,7}
#define PA_OUTPWR6_OID              {5,12,8}
#define CH6_VSWR_OID                {5,12,9}

#define RF_SWITCH7_OID              {5,13,1}
#define UP_GAIN_OFFSET7_OID         {5,13,2}
#define DN_GAIN_OFFSET7_OID         {5,13,3}
#define BAND_WIDE7_OID              {5,13,4}
#define UP_CENTER_FRE7_OID          {5,13,5}
#define DN_CENTER_FRE7_OID          {5,13,6}
#define UP_INPUT7_POWER_OID         {5,13,7}
#define PA_OUTPWR7_OID              {5,13,8}
#define CH7_VSWR_OID                {5,13,9}

#define RF_SWITCH8_OID              {5,14,1}
#define UP_GAIN_OFFSET8_OID         {5,14,2}
#define DN_GAIN_OFFSET8_OID         {5,14,3}
#define BAND_WIDE8_OID              {5,14,4}
#define UP_CENTER_FRE8_OID          {5,14,5}
#define DN_CENTER_FRE8_OID          {5,14,6}
#define UP_INPUT8_POWER_OID         {5,14,7}
#define PA_OUTPWR8_OID              {5,14,8}
#define CH8_VSWR_OID                {5,14,9}

/*N3RU*/
#define CH5_BASEBAND_UL_INPUT_POWER_OID  {5,11,10}
#define CH6_BASEBAND_UL_INPUT_POWER_OID  {5,12,10}
#define CH7_BASEBAND_UL_INPUT_POWER_OID  {5,13,10}
#define CH8_BASEBAND_UL_INPUT_POWER_OID  {5,14,10}

#define CH5_BASEBAND_DL_OUTPUT_POWER_OID {5,11,11}
#define CH6_BASEBAND_DL_OUTPUT_POWER_OID {5,12,11}
#define CH7_BASEBAND_DL_OUTPUT_POWER_OID {5,13,11}
#define CH8_BASEBAND_DL_OUTPUT_POWER_OID {5,14,11}



#define N2RU_DEV_TEMP_MAX_OID            {5,15,1}
#define N2RU_ROUTE_ADD_OID               {5,15,2}
#define FAN_WORK_MODEL_OID               {5,15,3}

// ------ for ru ------
#define CH1_UP_BASE_BAND_INPUT_POWER_OID            {5,16,1}
#define CH2_UP_BASE_BAND_INPUT_POWER_OID            {5,16,2}
#define CH3_UP_BASE_BAND_INPUT_POWER_OID            {5,16,3}
#define CH4_UP_BASE_BAND_INPUT_POWER_OID            {5,16,4}
#define CH5_UP_BASE_BAND_INPUT_POWER_OID            {5,16,5}
#define CH6_UP_BASE_BAND_INPUT_POWER_OID            {5,16,6}
#define CH7_UP_BASE_BAND_INPUT_POWER_OID            {5,16,7}
#define CH8_UP_BASE_BAND_INPUT_POWER_OID            {5,16,8}

#define CH1_DOWN_BASE_BAND_OUTPUT_POWER_OID          {5,16,9}
#define CH2_DOWN_BASE_BAND_OUTPUT_POWER_OID          {5,16,10}
#define CH3_DOWN_BASE_BAND_OUTPUT_POWER_OID          {5,16,11}
#define CH4_DOWN_BASE_BAND_OUTPUT_POWER_OID          {5,16,12}
#define CH5_DOWN_BASE_BAND_OUTPUT_POWER_OID          {5,16,13}
#define CH6_DOWN_BASE_BAND_OUTPUT_POWER_OID          {5,16,14}
#define CH7_DOWN_BASE_BAND_OUTPUT_POWER_OID          {5,16,15}
#define CH8_DOWN_BASE_BAND_OUTPUT_POWER_OID          {5,16,16}
// ------ end ------

// ------ for au ------
#define CH1_DOWN_BASE_BAND_INPUT_POWER_OID         {5,16,21}
#define CH2_DOWN_BASE_BAND_INPUT_POWER_OID         {5,16,22}
#define CH3_DOWN_BASE_BAND_INPUT_POWER_OID         {5,16,23}
#define CH4_DOWN_BASE_BAND_INPUT_POWER_OID         {5,16,24}

#define CH1_UP_BASE_BAND_OUTPUT_POWER_OID          {5,16,25}
#define CH2_UP_BASE_BAND_OUTPUT_POWER_OID          {5,16,26}
#define CH3_UP_BASE_BAND_OUTPUT_POWER_OID          {5,16,27}
#define CH4_UP_BASE_BAND_OUTPUT_POWER_OID          {5,16,28}
// ------ end ------

//----FRU--------
/* bti pa info 1 */
#define BTI_PA_1_MODEL_OID                         {5,17,1}
#define BTI_PA_1_SN_OID                            {5,17,2}
#define BTI_PA_1_FIRM_VER_OID                      {5,17,3}
#define BTI_PA_1_BAND_OID                          {5,17,4}
#define BTI_PA_1_VSWR_OID                          {5,17,5}
#define BTI_PA_1_DL_OUTPUT_POWER_OID               {5,17,6}
#define BTI_PA_1_TEMP_OID                          {5,17,7}
#define BTI_PA_1_RUN_STATUS_OID                    {5,17,8}
#define BTI_PA_1_UL_ATT_OID                        {5,17,9}
#define BTI_PA_1_DL_ATT_OID                        {5,17,10}
#define BTI_PA_1_UL_MAX_ATT_OID                    {5,17,11}
#define BTI_PA_1_PA_RESET_OID                      {5,17,12}
#define BTI_PA_1_PA_DISABLE_OID                    {5,17,13}
#define BTI_PA_1_DL_UNDER_POWER_ALM_EN_OID         {5,17,14}
#define BTI_PA_1_DL_UNDER_POWER_ALM_THR_OID        {5,17,15}
#define BTI_PA_1_UL_INPUT_POWER_OID                {5,17,16}
#define BTI_PA_1_UL_PILOT_FREQ_OID                 {5,17,17}
#define BTI_PA_1_UL_PILOT_TONE_EN_OID              {5,17,18}

/* bti pa info 2 */
#define BTI_PA_2_MODEL_OID                         {5,18,1}
#define BTI_PA_2_SN_OID                            {5,18,2}
#define BTI_PA_2_FIRM_VER_OID                      {5,18,3}
#define BTI_PA_2_BAND_OID                          {5,18,4}
#define BTI_PA_2_VSWR_OID                          {5,18,5}
#define BTI_PA_2_DL_OUTPUT_POWER_OID               {5,18,6}
#define BTI_PA_2_TEMP_OID                          {5,18,7}
#define BTI_PA_2_RUN_STATUS_OID                    {5,18,8}
#define BTI_PA_2_UL_ATT_OID                        {5,18,9}
#define BTI_PA_2_DL_ATT_OID                        {5,18,10}
#define BTI_PA_2_UL_MAX_ATT_OID                    {5,18,11}
#define BTI_PA_2_PA_RESET_OID                      {5,18,12}
#define BTI_PA_2_PA_DISABLE_OID                    {5,18,13}
#define BTI_PA_2_DL_UNDER_POWER_ALM_EN_OID         {5,18,14}
#define BTI_PA_2_DL_UNDER_POWER_ALM_THR_OID        {5,18,15}
#define BTI_PA_2_UL_INPUT_POWER_OID                {5,18,16}
#define BTI_PA_2_UL_PILOT_FREQ_OID                 {5,18,17}
#define BTI_PA_2_UL_PILOT_TONE_EN_OID              {5,18,18}

/* bti pa info 3 */
#define BTI_PA_3_MODEL_OID                         {5,19,1}
#define BTI_PA_3_SN_OID                            {5,19,2}
#define BTI_PA_3_FIRM_VER_OID                      {5,19,3}
#define BTI_PA_3_BAND_OID                          {5,19,4}
#define BTI_PA_3_VSWR_OID                          {5,19,5}
#define BTI_PA_3_DL_OUTPUT_POWER_OID               {5,19,6}
#define BTI_PA_3_TEMP_OID                          {5,19,7}
#define BTI_PA_3_RUN_STATUS_OID                    {5,19,8}
#define BTI_PA_3_UL_ATT_OID                        {5,19,9}
#define BTI_PA_3_DL_ATT_OID                        {5,19,10}
#define BTI_PA_3_UL_MAX_ATT_OID                    {5,19,11}
#define BTI_PA_3_PA_RESET_OID                      {5,19,12}
#define BTI_PA_3_PA_DISABLE_OID                    {5,19,13}
#define BTI_PA_3_DL_UNDER_POWER_ALM_EN_OID         {5,19,14}
#define BTI_PA_3_DL_UNDER_POWER_ALM_THR_OID        {5,19,15}
#define BTI_PA_3_UL_INPUT_POWER_OID                {5,19,16}
#define BTI_PA_3_UL_PILOT_FREQ_OID                 {5,19,17}
#define BTI_PA_3_UL_PILOT_TONE_EN_OID              {5,19,18}
                                                       
/* bti pa info 4 */
#define BTI_PA_4_MODEL_OID                         {5,20,1}
#define BTI_PA_4_SN_OID                            {5,20,2}
#define BTI_PA_4_FIRM_VER_OID                      {5,20,3}
#define BTI_PA_4_BAND_OID                          {5,20,4}
#define BTI_PA_4_VSWR_OID                          {5,20,5}
#define BTI_PA_4_DL_OUTPUT_POWER_OID               {5,20,6}
#define BTI_PA_4_TEMP_OID                          {5,20,7}
#define BTI_PA_4_RUN_STATUS_OID                    {5,20,8}
#define BTI_PA_4_UL_ATT_OID                        {5,20,9}
#define BTI_PA_4_DL_ATT_OID                        {5,20,10}
#define BTI_PA_4_UL_MAX_ATT_OID                    {5,20,11}
#define BTI_PA_4_PA_RESET_OID                      {5,20,12}
#define BTI_PA_4_PA_DISABLE_OID                    {5,20,13}
#define BTI_PA_4_DL_UNDER_POWER_ALM_EN_OID         {5,20,14}
#define BTI_PA_4_DL_UNDER_POWER_ALM_THR_OID        {5,20,15}
#define BTI_PA_4_UL_INPUT_POWER_OID                {5,20,16}
#define BTI_PA_4_UL_PILOT_FREQ_OID                 {5,20,17}
#define BTI_PA_4_UL_PILOT_TONE_EN_OID              {5,20,18}

#define POLYMERIC1_BASEBAND_UL_INPUT_POWER_OID	   {5,21,1}
#define POLYMERIC2_BASEBAND_UL_INPUT_POWER_OID	   {5,21,2}

#define WEBOMT_LOGOUT               {4,1,26}
#define A2_DEBUG_PORT              	{4,1,27}
#define SLAVE_DEBUG_PORT            {4,1,28}
#define SLAVE_WIFI_SWTICH           {4,1,29}


#define LIGHT_SERIAL_NUM_1			{4,6,1}
#define LIGHT_SERIAL_NUM_2			{4,6,2}
#define LIGHT_SERIAL_NUM_3			{4,6,3}
#define LIGHT_SERIAL_NUM_4			{4,6,4}
#define LIGHT_SERIAL_NUM_5			{4,6,5}
#define LIGHT_SERIAL_NUM_6			{4,6,6}
#define LIGHT_SERIAL_NUM_7			{4,6,7}
#define LIGHT_SERIAL_NUM_8			{4,6,8}
#define LIGHT_SERIAL_NUM_9			{4,6,9}
#define LIGHT_SERIAL_NUM_10			{4,6,10}
#define LIGHT_SERIAL_NUM_11			{4,6,11}
#define LIGHT_SERIAL_NUM_12			{4,6,12}
#define LIGHT_SERIAL_NUM_13			{4,6,13}
#define LIGHT_SERIAL_NUM_14			{4,6,14}
#define LIGHT_SERIAL_NUM_IN			{4,6,15} 
#define LIGHT_SERIAL_NUM_OUT		{4,6,16} 

#define LIGHT_1_Vender_Name			{4,6,17}
#define LIGHT_2_Vender_Name			{4,6,18}
#define LIGHT_3_Vender_Name			{4,6,19}
#define LIGHT_4_Vender_Name			{4,6,20}
#define LIGHT_5_Vender_Name			{4,6,21}
#define LIGHT_6_Vender_Name			{4,6,22}
#define LIGHT_7_Vender_Name			{4,6,23}
#define LIGHT_8_Vender_Name			{4,6,24}
#define LIGHT_9_Vender_Name			{4,6,25}
#define LIGHT_10_Vender_Name		{4,6,26}
#define LIGHT_11_Vender_Name		{4,6,27}
#define LIGHT_12_Vender_Name		{4,6,28}
#define LIGHT_13_Vender_Name		{4,6,29}
#define LIGHT_14_Vender_Name		{4,6,30}
#define LIGHT_IN_Vender_Name		{4,6,31}
#define LIGHT_OUT_Vender_Name		{4,6,32}

#define OP_1_TX_power				{4,6,33}
#define OP_2_TX_power				{4,6,34}
#define OP_3_TX_power				{4,6,35}
#define OP_4_TX_power				{4,6,36}
#define OP_5_TX_power				{4,6,37}
#define OP_6_TX_power				{4,6,38}
#define OP_7_TX_power				{4,6,39}
#define OP_8_TX_power				{4,6,40}
#define OP_9_TX_power				{4,6,41}
#define OP_10_TX_power				{4,6,42}
#define OP_11_TX_power				{4,6,43}
#define OP_12_TX_power				{4,6,44}
#define OP_13_TX_power				{4,6,45}
#define OP_14_TX_power				{4,6,46}
#define OP_slave_TX_power			{4,6,47}
#define OP_master_TX_power			{4,6,48}


#define OP_1_RX_power				{4,6,49}
#define OP_2_RX_power				{4,6,50}
#define OP_3_RX_power				{4,6,51}
#define OP_4_RX_power				{4,6,52}
#define OP_5_RX_power				{4,6,53}
#define OP_6_RX_power				{4,6,54}
#define OP_7_RX_power				{4,6,55}
#define OP_8_RX_power				{4,6,56}
#define OP_9_RX_power				{4,6,57}
#define OP_10_RX_power				{4,6,58}
#define OP_11_RX_power				{4,6,59}
#define OP_12_RX_power				{4,6,60}
#define OP_13_RX_power				{4,6,61}
#define OP_14_RX_power				{4,6,62}
#define OP_slave_RX_power			{4,6,63}
#define OP_master_RX_power			{4,6,64}

#define NTP_SWITCH      			{4,7,1}
#define NTP_UPDATE_INTERVAL      	{4,7,2}
#define NTP_TIME_ZONE				{4,7,3}
#define NTP_SERVER1_IP				{4,7,4}
#define NTP_SERVER2_IP				{4,7,5}

/***********************************************/

/***********************通道配置****************/

/*Master-AU时延*/
#define DELAY_TYPE                  {BAND_CONF_ROOT,1,1}
#define DELAY_ADJ_VAL               {BAND_CONF_ROOT,1,2}
#define DELAY_MEA_VAL               {BAND_CONF_ROOT,1,3}
#define DELAY_CONFIRM               {BAND_CONF_ROOT,1,4}
#define DELAY_CUS_VAL               {BAND_CONF_ROOT,1,5}

#define CARRIER_CONFIG_CONTROL      {BAND_CONF_ROOT,1,6}

#define TDD_SWITCH_5G_OID          {BAND_CONF_ROOT,2,1}
#define SYNC_STATUS_5G_OID         {BAND_CONF_ROOT,2,2}

#define MODULE1_5G_OID             {BAND_CONF_ROOT,2,3}
#define MODULE2_5G_OID             {BAND_CONF_ROOT,2,4}
#define MODULE3_5G_OID             {BAND_CONF_ROOT,2,5}
#define MODULE4_5G_OID             {BAND_CONF_ROOT,2,6}
#define SSB_TYPE_5G_OID            {BAND_CONF_ROOT,2,7}
/*XPG*/
#define SLOT_CONFIG_5G			    {BAND_CONF_ROOT,2,8}
#define SLOT_SYMBOLS_FORMAT1_5G		{BAND_CONF_ROOT,2,9}
#define SLOT_SYMBOLS_FORMAT2_5G		{BAND_CONF_ROOT,2,10}
#define SSB_ARFCN_5G_OID		    {BAND_CONF_ROOT,2,11}




#define TDD_SWITCH_4G_OID          {BAND_CONF_ROOT,2,13}
#define SYNC_STATUS_4G_OID         {BAND_CONF_ROOT,2,14}
#define MODULE1_4G_OID             {BAND_CONF_ROOT,2,15}
#define MODULE2_4G_OID             {BAND_CONF_ROOT,2,16}
#define MODULE3_4G_OID             {BAND_CONF_ROOT,2,17}
#define MODULE4_4G_OID             {BAND_CONF_ROOT,2,18}
#define UL_DL_CONFG_4G_OID         {BAND_CONF_ROOT,2,19}
#define SPCE_SUBFRAME_CONFG_4G_OID {BAND_CONF_ROOT,2,20}

#define BAND_CONF_MDL_1             (3)
#define BAND1_DATA_VALID_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_1,1}
#define BAND1_ULFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_1,2}
#define BAND1_ULFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_1,3}
#define BAND1_DLFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_1,4}
#define BAND1_DLFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_1,5}

#define BAND1_SIG1_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_1,6}
#define BAND1_SIG1_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_1,7}
#define BAND1_SIG1_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_1,8}
#define BAND1_SIG1_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_1,9}

#define BAND1_SIG2_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_1,10}
#define BAND1_SIG2_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_1,11}
#define BAND1_SIG2_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_1,12}
#define BAND1_SIG2_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_1,13}

#define BAND1_SIG3_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_1,14}
#define BAND1_SIG3_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_1,15}
#define BAND1_SIG3_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_1,16}
#define BAND1_SIG3_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_1,17}
#define CHL_1_SIGNAL_TRANS_BW_OID   {BAND_CONF_ROOT,BAND_CONF_MDL_1,18}

#define BAND_CONF_MDL_2             (4)
#define BAND2_DATA_VALID_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_2,1}
#define BAND2_ULFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_2,2}
#define BAND2_ULFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_2,3}
#define BAND2_DLFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_2,4}
#define BAND2_DLFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_2,5}

#define BAND2_SIG1_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_2,6}
#define BAND2_SIG1_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_2,7}
#define BAND2_SIG1_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_2,8}
#define BAND2_SIG1_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_2,9}

#define BAND2_SIG2_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_2,10}
#define BAND2_SIG2_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_2,11}
#define BAND2_SIG2_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_2,12}
#define BAND2_SIG2_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_2,13}

#define BAND2_SIG3_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_2,14}
#define BAND2_SIG3_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_2,15}
#define BAND2_SIG3_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_2,16}
#define BAND2_SIG3_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_2,17}
#define CHL_2_SIGNAL_TRANS_BW_OID   {BAND_CONF_ROOT,BAND_CONF_MDL_2,18}

#define BAND_CONF_MDL_3             (5)
#define BAND3_DATA_VALID_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_3,1}
#define BAND3_ULFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_3,2}
#define BAND3_ULFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_3,3}
#define BAND3_DLFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_3,4}
#define BAND3_DLFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_3,5}

#define BAND3_SIG1_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_3,6}
#define BAND3_SIG1_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_3,7}
#define BAND3_SIG1_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_3,8}
#define BAND3_SIG1_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_3,9}

#define BAND3_SIG2_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_3,10}
#define BAND3_SIG2_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_3,11}
#define BAND3_SIG2_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_3,12}
#define BAND3_SIG2_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_3,13}

#define BAND3_SIG3_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_3,14}
#define BAND3_SIG3_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_3,15}
#define BAND3_SIG3_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_3,16}
#define BAND3_SIG3_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_3,17}
#define CHL_3_SIGNAL_TRANS_BW_OID   {BAND_CONF_ROOT,BAND_CONF_MDL_3,18}

#define BAND_CONF_MDL_4             (6)
#define BAND4_DATA_VALID_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_4,1}
#define BAND4_ULFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_4,2}
#define BAND4_ULFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_4,3}
#define BAND4_DLFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_4,4}
#define BAND4_DLFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_4,5}

#define BAND4_SIG1_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_4,6}
#define BAND4_SIG1_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_4,7}
#define BAND4_SIG1_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_4,8}
#define BAND4_SIG1_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_4,9}

#define BAND4_SIG2_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_4,10}
#define BAND4_SIG2_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_4,11}
#define BAND4_SIG2_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_4,12}
#define BAND4_SIG2_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_4,13}

#define BAND4_SIG3_BAND_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_4,14}
#define BAND4_SIG3_TYPE_OID         {BAND_CONF_ROOT,BAND_CONF_MDL_4,15}
#define BAND4_SIG3_ULFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_4,16}
#define BAND4_SIG3_DLFREQ_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_4,17}
#define CHL_4_SIGNAL_TRANS_BW_OID   {BAND_CONF_ROOT,BAND_CONF_MDL_4,18}

/*N3RU*/
#define BAND_CONF_MDL_5             (9)
#define BAND5_ULFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_5,5}
#define BAND5_ULFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_5,6}
#define BAND5_DLFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_5,7}
#define BAND5_DLFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_5,8}
/*N3RU*/
#define BAND_CONF_MDL_6             (10)
#define BAND6_ULFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_6,5}
#define BAND6_ULFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_6,6}
#define BAND6_DLFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_6,7}
#define BAND6_DLFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_6,8}

/*N3RU*/
#define BAND_CONF_MDL_7             (11)
#define BAND7_ULFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_7,5}
#define BAND7_ULFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_7,6}
#define BAND7_DLFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_7,7}
#define BAND7_DLFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_7,8}

/*N3RU*/
#define BAND_CONF_MDL_8             (12)
#define BAND8_ULFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_8,5}
#define BAND8_ULFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_8,6}
#define BAND8_DLFREQ_LOW_OID        {BAND_CONF_ROOT,BAND_CONF_MDL_8,7}
#define BAND8_DLFREQ_HIGH_OID       {BAND_CONF_ROOT,BAND_CONF_MDL_8,8}


#define CHL_1_Transmission_Allocation   					{BAND_CONF_ROOT,BAND_CONF_MDL_1,19}
#define CHL_1_Carrier1_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_1,20}
#define CHL_1_Carrier1_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_1,21}
#define CHL_1_Carrier1_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_1,22}
#define CHL_1_Carrier1_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_1,23}
#define CHL_1_Carrier1_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_1,24}
#define CHL_1_Carrier1_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_1,25}
#define CHL_1_Carrier1_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_1,26}

#define CHL_1_Carrier2_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_1,27}
#define CHL_1_Carrier2_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_1,28}
#define CHL_1_Carrier2_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_1,29}
#define CHL_1_Carrier2_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_1,30}
#define CHL_1_Carrier2_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_1,31}
#define CHL_1_Carrier2_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_1,32}
#define CHL_1_Carrier2_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_1,33}

#define CHL_1_Carrier3_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_1,34}
#define CHL_1_Carrier3_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_1,35}
#define CHL_1_Carrier3_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_1,36}
#define CHL_1_Carrier3_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_1,37}
#define CHL_1_Carrier3_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_1,38}
#define CHL_1_Carrier3_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_1,39}
#define CHL_1_Carrier3_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_1,40}

#define CHL_1_Carrier4_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_1,41}
#define CHL_1_Carrier4_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_1,42}
#define CHL_1_Carrier4_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_1,43}
#define CHL_1_Carrier4_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_1,44}
#define CHL_1_Carrier4_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_1,45}
#define CHL_1_Carrier4_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_1,46}
#define CHL_1_Carrier4_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_1,47}

#define CHL_2_Transmission_Allocation   					{BAND_CONF_ROOT,BAND_CONF_MDL_2,19}
#define CHL_2_Carrier1_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_2,20}
#define CHL_2_Carrier1_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_2,21}
#define CHL_2_Carrier1_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_2,22}
#define CHL_2_Carrier1_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_2,23}
#define CHL_2_Carrier1_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_2,24}
#define CHL_2_Carrier1_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_2,25}
#define CHL_2_Carrier1_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_2,26}

#define CHL_2_Carrier2_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_2,27}
#define CHL_2_Carrier2_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_2,28}
#define CHL_2_Carrier2_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_2,29}
#define CHL_2_Carrier2_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_2,30}
#define CHL_2_Carrier2_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_2,31}
#define CHL_2_Carrier2_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_2,32}
#define CHL_2_Carrier2_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_2,33}

#define CHL_2_Carrier3_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_2,34}
#define CHL_2_Carrier3_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_2,35}
#define CHL_2_Carrier3_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_2,36}
#define CHL_2_Carrier3_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_2,37}
#define CHL_2_Carrier3_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_2,38}
#define CHL_2_Carrier3_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_2,39}
#define CHL_2_Carrier3_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_2,40}

#define CHL_2_Carrier4_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_2,41}
#define CHL_2_Carrier4_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_2,42}
#define CHL_2_Carrier4_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_2,43}
#define CHL_2_Carrier4_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_2,44}
#define CHL_2_Carrier4_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_2,45}
#define CHL_2_Carrier4_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_2,46}
#define CHL_2_Carrier4_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_2,47}

#define CHL_3_Transmission_Allocation   					{BAND_CONF_ROOT,BAND_CONF_MDL_3,19}
#define CHL_3_Carrier1_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_3,20}
#define CHL_3_Carrier1_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_3,21}
#define CHL_3_Carrier1_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_3,22}
#define CHL_3_Carrier1_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_3,23}
#define CHL_3_Carrier1_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_3,24}
#define CHL_3_Carrier1_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_3,25}
#define CHL_3_Carrier1_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_3,26}

#define CHL_3_Carrier2_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_3,27}
#define CHL_3_Carrier2_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_3,28}
#define CHL_3_Carrier2_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_3,29}
#define CHL_3_Carrier2_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_3,30}
#define CHL_3_Carrier2_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_3,31}
#define CHL_3_Carrier2_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_3,32}
#define CHL_3_Carrier2_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_3,33}

#define CHL_3_Carrier3_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_3,34}
#define CHL_3_Carrier3_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_3,35}
#define CHL_3_Carrier3_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_3,36}
#define CHL_3_Carrier3_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_3,37}
#define CHL_3_Carrier3_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_3,38}
#define CHL_3_Carrier3_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_3,39}
#define CHL_3_Carrier3_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_3,40}

#define CHL_3_Carrier4_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_3,41}
#define CHL_3_Carrier4_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_3,42}
#define CHL_3_Carrier4_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_3,43}
#define CHL_3_Carrier4_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_3,44}
#define CHL_3_Carrier4_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_3,45}
#define CHL_3_Carrier4_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_3,46}
#define CHL_3_Carrier4_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_3,47}

#define CHL_4_Transmission_Allocation   					{BAND_CONF_ROOT,BAND_CONF_MDL_4,19}
#define CHL_4_Carrier1_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_4,20}
#define CHL_4_Carrier1_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_4,21}
#define CHL_4_Carrier1_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_4,22}
#define CHL_4_Carrier1_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_4,23}
#define CHL_4_Carrier1_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_4,24}
#define CHL_4_Carrier1_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_4,25}
#define CHL_4_Carrier1_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_4,26}

#define CHL_4_Carrier2_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_4,27}
#define CHL_4_Carrier2_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_4,28}
#define CHL_4_Carrier2_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_4,29}
#define CHL_4_Carrier2_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_4,30}
#define CHL_4_Carrier2_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_4,31}
#define CHL_4_Carrier2_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_4,32}
#define CHL_4_Carrier2_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_4,33}

#define CHL_4_Carrier3_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_4,34}
#define CHL_4_Carrier3_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_4,35}
#define CHL_4_Carrier3_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_4,36}
#define CHL_4_Carrier3_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_4,37}
#define CHL_4_Carrier3_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_4,38}
#define CHL_4_Carrier3_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_4,39}
#define CHL_4_Carrier3_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_4,40}

#define CHL_4_Carrier4_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_4,41}
#define CHL_4_Carrier4_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_4,42}
#define CHL_4_Carrier4_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_4,43}
#define CHL_4_Carrier4_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_4,44}
#define CHL_4_Carrier4_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_4,45}
#define CHL_4_Carrier4_Module_Check							{BAND_CONF_ROOT,BAND_CONF_MDL_4,46}
#define CHL_4_Carrier4_Check_Info							{BAND_CONF_ROOT,BAND_CONF_MDL_4,47}

#define CHL_5_Carrier1_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_5,20}
#define CHL_5_Carrier1_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_5,21}
#define CHL_5_Carrier1_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_5,22}
#define CHL_5_Carrier1_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_5,23}
#define CHL_5_Carrier1_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_5,24}
#define CHL_5_Carrier2_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_5,27}
#define CHL_5_Carrier2_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_5,28}
#define CHL_5_Carrier2_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_5,29}
#define CHL_5_Carrier2_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_5,30}
#define CHL_5_Carrier2_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_5,31}
#define CHL_5_Carrier3_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_5,34}
#define CHL_5_Carrier3_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_5,35}
#define CHL_5_Carrier3_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_5,36}
#define CHL_5_Carrier3_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_5,37}
#define CHL_5_Carrier3_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_5,38}
#define CHL_5_Carrier4_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_5,41}
#define CHL_5_Carrier4_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_5,42}
#define CHL_5_Carrier4_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_5,43}
#define CHL_5_Carrier4_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_5,44}
#define CHL_5_Carrier4_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_5,45}

#define CHL_6_Carrier1_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_6,20}
#define CHL_6_Carrier1_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_6,21}
#define CHL_6_Carrier1_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_6,22}
#define CHL_6_Carrier1_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_6,23}
#define CHL_6_Carrier1_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_6,24}
#define CHL_6_Carrier2_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_6,27}
#define CHL_6_Carrier2_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_6,28}
#define CHL_6_Carrier2_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_6,29}
#define CHL_6_Carrier2_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_6,30}
#define CHL_6_Carrier2_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_6,31}
#define CHL_6_Carrier3_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_6,34}
#define CHL_6_Carrier3_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_6,35}
#define CHL_6_Carrier3_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_6,36}
#define CHL_6_Carrier3_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_6,37}
#define CHL_6_Carrier3_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_6,38}
#define CHL_6_Carrier4_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_6,41}
#define CHL_6_Carrier4_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_6,42}
#define CHL_6_Carrier4_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_6,43}
#define CHL_6_Carrier4_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_6,44}
#define CHL_6_Carrier4_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_6,45}

#define CHL_7_Carrier1_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_7,20}
#define CHL_7_Carrier1_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_7,21}
#define CHL_7_Carrier1_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_7,22}
#define CHL_7_Carrier1_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_7,23}
#define CHL_7_Carrier1_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_7,24}
#define CHL_7_Carrier2_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_7,27}
#define CHL_7_Carrier2_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_7,28}
#define CHL_7_Carrier2_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_7,29}
#define CHL_7_Carrier2_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_7,30}
#define CHL_7_Carrier2_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_7,31}
#define CHL_7_Carrier3_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_7,34}
#define CHL_7_Carrier3_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_7,35}
#define CHL_7_Carrier3_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_7,36}
#define CHL_7_Carrier3_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_7,37}
#define CHL_7_Carrier3_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_7,38}
#define CHL_7_Carrier4_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_7,41}
#define CHL_7_Carrier4_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_7,42}
#define CHL_7_Carrier4_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_7,43}
#define CHL_7_Carrier4_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_7,44}
#define CHL_7_Carrier4_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_7,45}

#define CHL_8_Carrier1_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_8,20}
#define CHL_8_Carrier1_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_8,21}
#define CHL_8_Carrier1_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_8,22}
#define CHL_8_Carrier1_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_8,23}
#define CHL_8_Carrier1_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_8,24}
#define CHL_8_Carrier2_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_8,27}
#define CHL_8_Carrier2_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_8,28}
#define CHL_8_Carrier2_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_8,29}
#define CHL_8_Carrier2_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_8,30}
#define CHL_8_Carrier2_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_8,31}
#define CHL_8_Carrier3_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_8,34}
#define CHL_8_Carrier3_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_8,35}
#define CHL_8_Carrier3_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_8,36}
#define CHL_8_Carrier3_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_8,37}
#define CHL_8_Carrier3_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_8,38}
#define CHL_8_Carrier4_UL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_8,41}
#define CHL_8_Carrier4_DL_Centre_Frequency					{BAND_CONF_ROOT,BAND_CONF_MDL_8,42}
#define CHL_8_Carrier4_Digital_Signal_Bandwidth_Select		{BAND_CONF_ROOT,BAND_CONF_MDL_8,43}
#define CHL_8_Carrier4_Digital_Signal_Bandwidth				{BAND_CONF_ROOT,BAND_CONF_MDL_8,44}
#define CHL_8_Carrier4_Carrier_Switch						{BAND_CONF_ROOT,BAND_CONF_MDL_8,45}

/*AU*/
#define BAND_CONFIG_UPDATE_OID      {BAND_CONF_ROOT,7,1}
#define LOC_SIG_BAND                {BAND_CONF_ROOT,7,2}
#define LOC_TRS_BAND                {BAND_CONF_ROOT,7,3}
#define SYS_SIG_BAND                {BAND_CONF_ROOT,7,4}
#define SYS_TRS_BAND                {BAND_CONF_ROOT,7,5}
#define GE_MODE_SELECT              {BAND_CONF_ROOT,7,6}
#define SYS_TRS_OCCUPIED_IN_FIBER   {BAND_CONF_ROOT,7,7}
#define MODE_800M_ENABLE            {BAND_CONF_ROOT,7,8}

/*RU*/
#define CHANNEL_MAP1                {BAND_CONF_ROOT,8,1}
#define CHANNEL_MAP2                {BAND_CONF_ROOT,8,2}
#define CHANNEL_MAP3                {BAND_CONF_ROOT,8,3}
#define CHANNEL_MAP4                {BAND_CONF_ROOT,8,4}

/*A3*/
#define CARRIER_CONFIG_SWITCH       {BAND_CONF_ROOT,8,5}

/* SU-HP */
#define CHANNEL_MAP1_NEW            {BAND_CONF_ROOT,8,5}
#define CHANNEL_MAP2_NEW            {BAND_CONF_ROOT,8,6}

/* N2RU */
#define CHANNEL1_FREQ_RANGE         {BAND_CONF_ROOT,8,7}
#define CHANNEL1_MAP_V2             {BAND_CONF_ROOT,8,8}
#define CHANNEL2_FREQ_RANGE         {BAND_CONF_ROOT,8,9}
#define CHANNEL2_MAP_V2             {BAND_CONF_ROOT,8,10}
#define CHANNEL3_FREQ_RANGE         {BAND_CONF_ROOT,8,11}
#define CHANNEL3_MAP_V2             {BAND_CONF_ROOT,8,12}
#define CHANNEL4_FREQ_RANGE         {BAND_CONF_ROOT,8,13}
#define CHANNEL4_MAP_V2             {BAND_CONF_ROOT,8,14}
#define CHANNEL5_FREQ_RANGE         {BAND_CONF_ROOT,8,15}
#define CHANNEL5_MAP_V2             {BAND_CONF_ROOT,8,16}
#define CHANNEL6_FREQ_RANGE         {BAND_CONF_ROOT,8,17}
#define CHANNEL6_MAP_V2             {BAND_CONF_ROOT,8,18}
#define CHANNEL7_FREQ_RANGE         {BAND_CONF_ROOT,8,19}
#define CHANNEL7_MAP_V2             {BAND_CONF_ROOT,8,20}
#define CHANNEL8_FREQ_RANGE         {BAND_CONF_ROOT,8,21}
#define CHANNEL8_MAP_V2             {BAND_CONF_ROOT,8,22}
/*SU*/
#define CH1_TDD_CONFIG_ROOT	(9)
#define CH1_TDD_ENABLE_OID		{BAND_CONF_ROOT, CH1_TDD_CONFIG_ROOT, 1}
#define CH1_ULDL_CONFIG_OID		{BAND_CONF_ROOT, CH1_TDD_CONFIG_ROOT, 2}
#define CH1_SUBFRAMECONF_OID	{BAND_CONF_ROOT, CH1_TDD_CONFIG_ROOT, 3}
#define CH1_CP_TYPE_OID			{BAND_CONF_ROOT, CH1_TDD_CONFIG_ROOT, 4}

#define CH2_TDD_CONFIG_ROOT	(10)
#define CH2_TDD_ENABLE_OID		{BAND_CONF_ROOT, CH2_TDD_CONFIG_ROOT, 1}
#define CH2_ULDL_CONFIG_OID		{BAND_CONF_ROOT, CH2_TDD_CONFIG_ROOT, 2}
#define CH2_SUBFRAMECONF_OID	{BAND_CONF_ROOT, CH2_TDD_CONFIG_ROOT, 3}
#define CH2_CP_TYPE_OID			{BAND_CONF_ROOT, CH2_TDD_CONFIG_ROOT, 4}

#define CH3_TDD_CONFIG_ROOT	(11)
#define CH3_TDD_ENABLE_OID		{BAND_CONF_ROOT, CH3_TDD_CONFIG_ROOT, 1}
#define CH3_ULDL_CONFIG_OID		{BAND_CONF_ROOT, CH3_TDD_CONFIG_ROOT, 2}
#define CH3_SUBFRAMECONF_OID	{BAND_CONF_ROOT, CH3_TDD_CONFIG_ROOT, 3}
#define CH3_CP_TYPE_OID			{BAND_CONF_ROOT, CH3_TDD_CONFIG_ROOT, 4}

#define CH4_TDD_CONFIG_ROOT	(12)
#define CH4_TDD_ENABLE_OID		{BAND_CONF_ROOT, CH4_TDD_CONFIG_ROOT, 1}
#define CH4_ULDL_CONFIG_OID		{BAND_CONF_ROOT, CH4_TDD_CONFIG_ROOT, 2}
#define CH4_SUBFRAMECONF_OID	{BAND_CONF_ROOT, CH4_TDD_CONFIG_ROOT, 3}
#define CH4_CP_TYPE_OID			{BAND_CONF_ROOT, CH4_TDD_CONFIG_ROOT, 4}

/***********************************************************/
/*****************AU commbiner******************************/
#define POI_CONTROL_MODE                {POI_ROOT,1,1}
#define POI_ADJUST_INTER                {POI_ROOT,1,2}
#define POI_ATT_RESET                	{POI_ROOT,1,3}

#define CHAN1_POI_POWER_SWITCH			{POI_ROOT,2,1}
#define CHAN1PORT1_INPUT_POWER          {POI_ROOT,2,2}
#define CHAN1PORT2_INPUT_POWER          {POI_ROOT,2,3}
#define CHAN1PORT3_INPUT_POWER          {POI_ROOT,2,4}
#define CHAN1PORT4_INPUT_POWER          {POI_ROOT,2,5}
#define CHAN1PORT1_ATTENUATION          {POI_ROOT,2,6}
#define CHAN1PORT2_ATTENUATION          {POI_ROOT,2,7}
#define CHAN1PORT3_ATTENUATION          {POI_ROOT,2,8}
#define CHAN1PORT4_ATTENUATION          {POI_ROOT,2,9}
#define CHAN1PORT1_OFFSET               {POI_ROOT,2,10}
#define CHAN1PORT2_OFFSET               {POI_ROOT,2,11}
#define CHAN1PORT3_OFFSET               {POI_ROOT,2,12}
#define CHAN1PORT4_OFFSET               {POI_ROOT,2,13}
#define CHAN1PORT1_OPERATOR             {POI_ROOT,2,14}
#define CHAN1PORT2_OPERATOR             {POI_ROOT,2,15}
#define CHAN1PORT3_OPERATOR             {POI_ROOT,2,16}
#define CHAN1PORT4_OPERATOR             {POI_ROOT,2,17}
#define CHAN1_POI_SERIAL_NUM		    {POI_ROOT,2,18}
#define CHAN1PORT1_SERVICE              {POI_ROOT,2,19}
#define CHAN1PORT2_SERVICE              {POI_ROOT,2,20}
#define CHAN1PORT3_SERVICE              {POI_ROOT,2,21}
#define CHAN1PORT4_SERVICE              {POI_ROOT,2,22}


#define CHAN2_POI_POWER_SWITCH			{POI_ROOT,3,1}
#define CHAN2PORT1_INPUT_POWER          {POI_ROOT,3,2}
#define CHAN2PORT2_INPUT_POWER          {POI_ROOT,3,3}
#define CHAN2PORT3_INPUT_POWER          {POI_ROOT,3,4}
#define CHAN2PORT4_INPUT_POWER          {POI_ROOT,3,5}
#define CHAN2PORT1_ATTENUATION          {POI_ROOT,3,6}
#define CHAN2PORT2_ATTENUATION          {POI_ROOT,3,7}
#define CHAN2PORT3_ATTENUATION          {POI_ROOT,3,8}
#define CHAN2PORT4_ATTENUATION          {POI_ROOT,3,9}
#define CHAN2PORT1_OFFSET               {POI_ROOT,3,10}
#define CHAN2PORT2_OFFSET               {POI_ROOT,3,11}
#define CHAN2PORT3_OFFSET               {POI_ROOT,3,12}
#define CHAN2PORT4_OFFSET               {POI_ROOT,3,13}
#define CHAN2PORT1_OPERATOR             {POI_ROOT,3,14}
#define CHAN2PORT2_OPERATOR             {POI_ROOT,3,15}
#define CHAN2PORT3_OPERATOR             {POI_ROOT,3,16}
#define CHAN2PORT4_OPERATOR             {POI_ROOT,3,17}
#define CHAN2_POI_SERIAL_NUM		    {POI_ROOT,3,18}
#define CHAN2PORT1_SERVICE              {POI_ROOT,3,19}
#define CHAN2PORT2_SERVICE              {POI_ROOT,3,20}
#define CHAN2PORT3_SERVICE              {POI_ROOT,3,21}
#define CHAN2PORT4_SERVICE              {POI_ROOT,3,22}


#define CHAN3_POI_POWER_SWITCH			{POI_ROOT,4,1}
#define CHAN3PORT1_INPUT_POWER          {POI_ROOT,4,2}
#define CHAN3PORT2_INPUT_POWER          {POI_ROOT,4,3}
#define CHAN3PORT3_INPUT_POWER          {POI_ROOT,4,4}
#define CHAN3PORT4_INPUT_POWER          {POI_ROOT,4,5}
#define CHAN3PORT1_ATTENUATION          {POI_ROOT,4,6}
#define CHAN3PORT2_ATTENUATION          {POI_ROOT,4,7}
#define CHAN3PORT3_ATTENUATION          {POI_ROOT,4,8}
#define CHAN3PORT4_ATTENUATION          {POI_ROOT,4,9}
#define CHAN3PORT1_OFFSET               {POI_ROOT,4,10}
#define CHAN3PORT2_OFFSET               {POI_ROOT,4,11}
#define CHAN3PORT3_OFFSET               {POI_ROOT,4,12}
#define CHAN3PORT4_OFFSET               {POI_ROOT,4,13}
#define CHAN3PORT1_OPERATOR             {POI_ROOT,4,14}
#define CHAN3PORT2_OPERATOR             {POI_ROOT,4,15}
#define CHAN3PORT3_OPERATOR             {POI_ROOT,4,16}
#define CHAN3PORT4_OPERATOR             {POI_ROOT,4,17}
#define CHAN3_POI_SERIAL_NUM		    {POI_ROOT,4,18}
#define CHAN3PORT1_SERVICE              {POI_ROOT,4,19}
#define CHAN3PORT2_SERVICE              {POI_ROOT,4,20}
#define CHAN3PORT3_SERVICE              {POI_ROOT,4,21}
#define CHAN3PORT4_SERVICE              {POI_ROOT,4,22}


#define CHAN4_POI_POWER_SWITCH			{POI_ROOT,5,1}
#define CHAN4PORT1_INPUT_POWER          {POI_ROOT,5,2}
#define CHAN4PORT2_INPUT_POWER          {POI_ROOT,5,3}
#define CHAN4PORT3_INPUT_POWER          {POI_ROOT,5,4}
#define CHAN4PORT4_INPUT_POWER          {POI_ROOT,5,5}
#define CHAN4PORT1_ATTENUATION          {POI_ROOT,5,6}
#define CHAN4PORT2_ATTENUATION          {POI_ROOT,5,7}
#define CHAN4PORT3_ATTENUATION          {POI_ROOT,5,8}
#define CHAN4PORT4_ATTENUATION          {POI_ROOT,5,9}
#define CHAN4PORT1_OFFSET               {POI_ROOT,5,10}
#define CHAN4PORT2_OFFSET               {POI_ROOT,5,11}
#define CHAN4PORT3_OFFSET               {POI_ROOT,5,12}
#define CHAN4PORT4_OFFSET               {POI_ROOT,5,13}
#define CHAN4PORT1_OPERATOR             {POI_ROOT,5,14}
#define CHAN4PORT2_OPERATOR             {POI_ROOT,5,15}
#define CHAN4PORT3_OPERATOR             {POI_ROOT,5,16}
#define CHAN4PORT4_OPERATOR             {POI_ROOT,5,17}
#define CHAN4_POI_SERIAL_NUM		    {POI_ROOT,5,18}
#define CHAN4PORT1_SERVICE              {POI_ROOT,5,19}
#define CHAN4PORT2_SERVICE              {POI_ROOT,5,20}
#define CHAN4PORT3_SERVICE              {POI_ROOT,5,21}
#define CHAN4PORT4_SERVICE              {POI_ROOT,5,22}



#define ARU_PA_SERIAL_NUM_1            {7,3,21}
#define ARU_PA_SERIAL_NUM_2            {7,4,22}
#define ARU_PA_SERIAL_NUM_3            {7,5,23}
#define ARU_PA_SERIAL_NUM_4            {7,6,24}

#define PA_MDL1_SERIAL_NUM            {5,1,10}
#define PA_MDL2_SERIAL_NUM            {5,2,10}
#define PA_MDL3_SERIAL_NUM            {5,3,10}
#define PA_MDL4_SERIAL_NUM            {5,4,10}

/********************************************************/

/******************A3 Carrier Power**********************/
#define CARRIER_POWER_CONFIG_ROOT	(9)
#define CHL1_C1_UL_Att							{CARRIER_POWER_CONFIG_ROOT,1,1}
#define CHL1_C1_DL_Att							{CARRIER_POWER_CONFIG_ROOT,1,2}
#define CHL1_C1_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,1,3}
#define CHL1_C1_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,1,4}
#define CHL1_C2_UL_Att							{CARRIER_POWER_CONFIG_ROOT,1,5}
#define CHL1_C2_DL_Att							{CARRIER_POWER_CONFIG_ROOT,1,6}
#define CHL1_C2_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,1,7}
#define CHL1_C2_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,1,8}
#define CHL1_C3_UL_Att							{CARRIER_POWER_CONFIG_ROOT,1,9}
#define CHL1_C3_DL_Att							{CARRIER_POWER_CONFIG_ROOT,1,10}
#define CHL1_C3_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,1,11}
#define CHL1_C3_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,1,12}
#define CHL1_C4_UL_Att							{CARRIER_POWER_CONFIG_ROOT,1,13}
#define CHL1_C4_DL_Att							{CARRIER_POWER_CONFIG_ROOT,1,14}
#define CHL1_C4_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,1,15}
#define CHL1_C4_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,1,16}

#define CHL2_C1_UL_Att							{CARRIER_POWER_CONFIG_ROOT,2,1}
#define CHL2_C1_DL_Att							{CARRIER_POWER_CONFIG_ROOT,2,2}
#define CHL2_C1_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,2,3}
#define CHL2_C1_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,2,4}
#define CHL2_C2_UL_Att							{CARRIER_POWER_CONFIG_ROOT,2,5}
#define CHL2_C2_DL_Att							{CARRIER_POWER_CONFIG_ROOT,2,6}
#define CHL2_C2_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,2,7}
#define CHL2_C2_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,2,8}
#define CHL2_C3_UL_Att							{CARRIER_POWER_CONFIG_ROOT,2,9}
#define CHL2_C3_DL_Att							{CARRIER_POWER_CONFIG_ROOT,2,10}
#define CHL2_C3_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,2,11}
#define CHL2_C3_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,2,12}
#define CHL2_C4_UL_Att							{CARRIER_POWER_CONFIG_ROOT,2,13}
#define CHL2_C4_DL_Att							{CARRIER_POWER_CONFIG_ROOT,2,14}
#define CHL2_C4_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,2,15}
#define CHL2_C4_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,2,16}

#define CHL3_C1_UL_Att							{CARRIER_POWER_CONFIG_ROOT,3,1}
#define CHL3_C1_DL_Att							{CARRIER_POWER_CONFIG_ROOT,3,2}
#define CHL3_C1_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,3,3}
#define CHL3_C1_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,3,4}
#define CHL3_C2_UL_Att							{CARRIER_POWER_CONFIG_ROOT,3,5}
#define CHL3_C2_DL_Att							{CARRIER_POWER_CONFIG_ROOT,3,6}
#define CHL3_C2_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,3,7}
#define CHL3_C2_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,3,8}
#define CHL3_C3_UL_Att							{CARRIER_POWER_CONFIG_ROOT,3,9}
#define CHL3_C3_DL_Att							{CARRIER_POWER_CONFIG_ROOT,3,10}
#define CHL3_C3_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,3,11}
#define CHL3_C3_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,3,12}
#define CHL3_C4_UL_Att							{CARRIER_POWER_CONFIG_ROOT,3,13}
#define CHL3_C4_DL_Att							{CARRIER_POWER_CONFIG_ROOT,3,14}
#define CHL3_C4_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,3,15}
#define CHL3_C4_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,3,16}

#define CHL4_C1_UL_Att							{CARRIER_POWER_CONFIG_ROOT,4,1}
#define CHL4_C1_DL_Att							{CARRIER_POWER_CONFIG_ROOT,4,2}
#define CHL4_C1_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,4,3}
#define CHL4_C1_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,4,4}
#define CHL4_C2_UL_Att							{CARRIER_POWER_CONFIG_ROOT,4,5}
#define CHL4_C2_DL_Att							{CARRIER_POWER_CONFIG_ROOT,4,6}
#define CHL4_C2_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,4,7}
#define CHL4_C2_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,4,8}
#define CHL4_C3_UL_Att							{CARRIER_POWER_CONFIG_ROOT,4,9}
#define CHL4_C3_DL_Att							{CARRIER_POWER_CONFIG_ROOT,4,10}
#define CHL4_C3_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,4,11}
#define CHL4_C3_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,4,12}
#define CHL4_C4_UL_Att							{CARRIER_POWER_CONFIG_ROOT,4,13}
#define CHL4_C4_DL_Att							{CARRIER_POWER_CONFIG_ROOT,4,14}
#define CHL4_C4_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,4,15}
#define CHL4_C4_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,4,16}

#define CHL5_C1_UL_Att							{CARRIER_POWER_CONFIG_ROOT,5,1}
#define CHL5_C1_DL_Att							{CARRIER_POWER_CONFIG_ROOT,5,2}
#define CHL5_C1_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,5,3}
#define CHL5_C1_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,5,4}
#define CHL5_C2_UL_Att							{CARRIER_POWER_CONFIG_ROOT,5,5}
#define CHL5_C2_DL_Att							{CARRIER_POWER_CONFIG_ROOT,5,6}
#define CHL5_C2_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,5,7}
#define CHL5_C2_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,5,8}
#define CHL5_C3_UL_Att							{CARRIER_POWER_CONFIG_ROOT,5,9}
#define CHL5_C3_DL_Att							{CARRIER_POWER_CONFIG_ROOT,5,10}
#define CHL5_C3_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,5,11}
#define CHL5_C3_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,5,12}
#define CHL5_C4_UL_Att							{CARRIER_POWER_CONFIG_ROOT,5,13}
#define CHL5_C4_DL_Att							{CARRIER_POWER_CONFIG_ROOT,5,14}
#define CHL5_C4_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,5,15}
#define CHL5_C4_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,5,16}

#define CHL6_C1_UL_Att							{CARRIER_POWER_CONFIG_ROOT,6,1}
#define CHL6_C1_DL_Att							{CARRIER_POWER_CONFIG_ROOT,6,2}
#define CHL6_C1_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,6,3}
#define CHL6_C1_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,6,4}
#define CHL6_C2_UL_Att							{CARRIER_POWER_CONFIG_ROOT,6,5}
#define CHL6_C2_DL_Att							{CARRIER_POWER_CONFIG_ROOT,6,6}
#define CHL6_C2_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,6,7}
#define CHL6_C2_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,6,8}
#define CHL6_C3_UL_Att							{CARRIER_POWER_CONFIG_ROOT,6,9}
#define CHL6_C3_DL_Att							{CARRIER_POWER_CONFIG_ROOT,6,10}
#define CHL6_C3_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,6,11}
#define CHL6_C3_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,6,12}
#define CHL6_C4_UL_Att							{CARRIER_POWER_CONFIG_ROOT,6,13}
#define CHL6_C4_DL_Att							{CARRIER_POWER_CONFIG_ROOT,6,14}
#define CHL6_C4_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,6,15}
#define CHL6_C4_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,6,16}

#define CHL7_C1_UL_Att							{CARRIER_POWER_CONFIG_ROOT,7,1}
#define CHL7_C1_DL_Att							{CARRIER_POWER_CONFIG_ROOT,7,2}
#define CHL7_C1_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,7,3}
#define CHL7_C1_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,7,4}
#define CHL7_C2_UL_Att							{CARRIER_POWER_CONFIG_ROOT,7,5}
#define CHL7_C2_DL_Att							{CARRIER_POWER_CONFIG_ROOT,7,6}
#define CHL7_C2_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,7,7}
#define CHL7_C2_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,7,8}
#define CHL7_C3_UL_Att							{CARRIER_POWER_CONFIG_ROOT,7,9}
#define CHL7_C3_DL_Att							{CARRIER_POWER_CONFIG_ROOT,7,10}
#define CHL7_C3_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,7,11}
#define CHL7_C3_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,7,12}
#define CHL7_C4_UL_Att							{CARRIER_POWER_CONFIG_ROOT,7,13}
#define CHL7_C4_DL_Att							{CARRIER_POWER_CONFIG_ROOT,7,14}
#define CHL7_C4_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,7,15}
#define CHL7_C4_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,7,16}

#define CHL8_C1_UL_Att							{CARRIER_POWER_CONFIG_ROOT,8,1}
#define CHL8_C1_DL_Att							{CARRIER_POWER_CONFIG_ROOT,8,2}
#define CHL8_C1_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,8,3}
#define CHL8_C1_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,8,4}
#define CHL8_C2_UL_Att							{CARRIER_POWER_CONFIG_ROOT,8,5}
#define CHL8_C2_DL_Att							{CARRIER_POWER_CONFIG_ROOT,8,6}
#define CHL8_C2_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,8,7}
#define CHL8_C2_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,8,8}
#define CHL8_C3_UL_Att							{CARRIER_POWER_CONFIG_ROOT,8,9}
#define CHL8_C3_DL_Att							{CARRIER_POWER_CONFIG_ROOT,8,10}
#define CHL8_C3_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,8,11}
#define CHL8_C3_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,8,12}
#define CHL8_C4_UL_Att							{CARRIER_POWER_CONFIG_ROOT,8,13}
#define CHL8_C4_DL_Att							{CARRIER_POWER_CONFIG_ROOT,8,14}
#define CHL8_C4_UL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,8,15}
#define CHL8_C4_DL_Carrier_Power				{CARRIER_POWER_CONFIG_ROOT,8,16}

/********************************************************/


/******************RU module-port-gain*******************/
#define CHANNEL1_PORT1_ULGAIN        {GAIN_ROOT,1,1}
#define CHANNEL1_PORT2_ULGAIN        {GAIN_ROOT,1,2}
#define CHANNEL1_PORT3_ULGAIN        {GAIN_ROOT,1,3}
#define CHANNEL1_PORT4_ULGAIN        {GAIN_ROOT,1,4}
#define CHANNEL1_PORT1_DLGAIN        {GAIN_ROOT,1,5}
#define CHANNEL1_PORT2_DLGAIN        {GAIN_ROOT,1,6}
#define CHANNEL1_PORT3_DLGAIN        {GAIN_ROOT,1,7}
#define CHANNEL1_PORT4_DLGAIN        {GAIN_ROOT,1,8}

#define CHANNEL2_PORT1_ULGAIN        {GAIN_ROOT,2,1}
#define CHANNEL2_PORT2_ULGAIN        {GAIN_ROOT,2,2}
#define CHANNEL2_PORT3_ULGAIN        {GAIN_ROOT,2,3}
#define CHANNEL2_PORT4_ULGAIN        {GAIN_ROOT,2,4}
#define CHANNEL2_PORT1_DLGAIN        {GAIN_ROOT,2,5}
#define CHANNEL2_PORT2_DLGAIN        {GAIN_ROOT,2,6}
#define CHANNEL2_PORT3_DLGAIN        {GAIN_ROOT,2,7}
#define CHANNEL2_PORT4_DLGAIN        {GAIN_ROOT,2,8}

#define CHANNEL3_PORT1_ULGAIN        {GAIN_ROOT,3,1}
#define CHANNEL3_PORT2_ULGAIN        {GAIN_ROOT,3,2}
#define CHANNEL3_PORT3_ULGAIN        {GAIN_ROOT,3,3}
#define CHANNEL3_PORT4_ULGAIN        {GAIN_ROOT,3,4}
#define CHANNEL3_PORT1_DLGAIN        {GAIN_ROOT,3,5}
#define CHANNEL3_PORT2_DLGAIN        {GAIN_ROOT,3,6}
#define CHANNEL3_PORT3_DLGAIN        {GAIN_ROOT,3,7}
#define CHANNEL3_PORT4_DLGAIN        {GAIN_ROOT,3,8}

#define CHANNEL4_PORT1_ULGAIN        {GAIN_ROOT,4,1}
#define CHANNEL4_PORT2_ULGAIN        {GAIN_ROOT,4,2}
#define CHANNEL4_PORT3_ULGAIN        {GAIN_ROOT,4,3}
#define CHANNEL4_PORT4_ULGAIN        {GAIN_ROOT,4,4}
#define CHANNEL4_PORT1_DLGAIN        {GAIN_ROOT,4,5}
#define CHANNEL4_PORT2_DLGAIN        {GAIN_ROOT,4,6}
#define CHANNEL4_PORT3_DLGAIN        {GAIN_ROOT,4,7}
#define CHANNEL4_PORT4_DLGAIN        {GAIN_ROOT,4,8}
/* N2RU */
#define CHANNEL5_PORT1_ULGAIN        {GAIN_ROOT,5,1}
#define CHANNEL5_PORT2_ULGAIN        {GAIN_ROOT,5,2}
#define CHANNEL5_PORT3_ULGAIN        {GAIN_ROOT,5,3}
#define CHANNEL5_PORT4_ULGAIN        {GAIN_ROOT,5,4}
#define CHANNEL5_PORT1_DLGAIN        {GAIN_ROOT,5,5}
#define CHANNEL5_PORT2_DLGAIN        {GAIN_ROOT,5,6}
#define CHANNEL5_PORT3_DLGAIN        {GAIN_ROOT,5,7}
#define CHANNEL5_PORT4_DLGAIN        {GAIN_ROOT,5,8}

#define CHANNEL6_PORT1_ULGAIN        {GAIN_ROOT,6,1}
#define CHANNEL6_PORT2_ULGAIN        {GAIN_ROOT,6,2}
#define CHANNEL6_PORT3_ULGAIN        {GAIN_ROOT,6,3}
#define CHANNEL6_PORT4_ULGAIN        {GAIN_ROOT,6,4}
#define CHANNEL6_PORT1_DLGAIN        {GAIN_ROOT,6,5}
#define CHANNEL6_PORT2_DLGAIN        {GAIN_ROOT,6,6}
#define CHANNEL6_PORT3_DLGAIN        {GAIN_ROOT,6,7}
#define CHANNEL6_PORT4_DLGAIN        {GAIN_ROOT,6,8}

#define CHANNEL7_PORT1_ULGAIN        {GAIN_ROOT,7,1}
#define CHANNEL7_PORT2_ULGAIN        {GAIN_ROOT,7,2}
#define CHANNEL7_PORT3_ULGAIN        {GAIN_ROOT,7,3}
#define CHANNEL7_PORT4_ULGAIN        {GAIN_ROOT,7,4}
#define CHANNEL7_PORT1_DLGAIN        {GAIN_ROOT,7,5}
#define CHANNEL7_PORT2_DLGAIN        {GAIN_ROOT,7,6}
#define CHANNEL7_PORT3_DLGAIN        {GAIN_ROOT,7,7}
#define CHANNEL7_PORT4_DLGAIN        {GAIN_ROOT,7,8}

#define CHANNEL8_PORT1_ULGAIN        {GAIN_ROOT,8,1}
#define CHANNEL8_PORT2_ULGAIN        {GAIN_ROOT,8,2}
#define CHANNEL8_PORT3_ULGAIN        {GAIN_ROOT,8,3}
#define CHANNEL8_PORT4_ULGAIN        {GAIN_ROOT,8,4}
#define CHANNEL8_PORT1_DLGAIN        {GAIN_ROOT,8,5}
#define CHANNEL8_PORT2_DLGAIN        {GAIN_ROOT,8,6}
#define CHANNEL8_PORT3_DLGAIN        {GAIN_ROOT,8,7}
#define CHANNEL8_PORT4_DLGAIN        {GAIN_ROOT,8,8}
/******************************************************/

/*******************RU 容量调度节点***********************/
#define CAPACITY_ALLOCATION_SW_OID          {CAPACITY_ROOT,1,1}
#define AU_1_INFO_UL_OID                    {CAPACITY_ROOT,2,1}
#define AU_1_INFO_DL_OID                    {CAPACITY_ROOT,2,2}
#define AU_2_INFO_UL_OID                    {CAPACITY_ROOT,2,3}
#define AU_2_INFO_DL_OID                    {CAPACITY_ROOT,2,4}
#define AU_3_INFO_UL_OID                    {CAPACITY_ROOT,2,5}
#define AU_3_INFO_DL_OID                    {CAPACITY_ROOT,2,6}
#define AU_4_INFO_UL_OID                    {CAPACITY_ROOT,2,7}
#define AU_4_INFO_DL_OID                    {CAPACITY_ROOT,2,8}
#define SAU1_1_INFO_UL_OID                  {CAPACITY_ROOT,2,9}
#define SAU1_1_INFO_DL_OID                  {CAPACITY_ROOT,2,10}
#define SAU1_2_INFO_UL_OID                  {CAPACITY_ROOT,2,11}
#define SAU1_2_INFO_DL_OID                  {CAPACITY_ROOT,2,12}
#define SAU1_3_INFO_UL_OID                  {CAPACITY_ROOT,2,13}
#define SAU1_3_INFO_DL_OID                  {CAPACITY_ROOT,2,14}
#define SAU1_4_INFO_UL_OID                  {CAPACITY_ROOT,2,15}
#define SAU1_4_INFO_DL_OID                  {CAPACITY_ROOT,2,16}
#define SAU2_1_INFO_UL_OID                  {CAPACITY_ROOT,2,17}
#define SAU2_1_INFO_DL_OID                  {CAPACITY_ROOT,2,18}
#define SAU2_2_INFO_UL_OID                  {CAPACITY_ROOT,2,19}
#define SAU2_2_INFO_DL_OID                  {CAPACITY_ROOT,2,20}
#define SAU2_3_INFO_UL_OID                  {CAPACITY_ROOT,2,21}
#define SAU2_3_INFO_DL_OID                  {CAPACITY_ROOT,2,22}
#define SAU2_4_INFO_UL_OID                  {CAPACITY_ROOT,2,23}
#define SAU2_4_INFO_DL_OID                  {CAPACITY_ROOT,2,24}

#define CH1_INFO_ULDL_OID                   {CAPACITY_ROOT,3,1}
#define CH2_INFO_ULDL_OID                   {CAPACITY_ROOT,3,2}
#define CH3_INFO_ULDL_OID                   {CAPACITY_ROOT,3,3}
#define CH4_INFO_ULDL_OID                   {CAPACITY_ROOT,3,4}

#define HRU_RF_GROUP1_OID                   {CAPACITY_ROOT,4,1}
#define HRU_RF_GROUP2_OID                   {CAPACITY_ROOT,4,2}
#define HRU_RF_GROUP3_OID                   {CAPACITY_ROOT,4,3}

#define MRU_RF_GROUP1_OID                   {CAPACITY_ROOT,4,5}
#define MRU_RF_GROUP2_OID                   {CAPACITY_ROOT,4,6}
#define MRU_RF_GROUP3_OID                   {CAPACITY_ROOT,4,7}

#define RU_RF_GROUP1_OID                    {CAPACITY_ROOT,4,8}
#define RU_RF_GROUP2_OID                    {CAPACITY_ROOT,4,9}
#define RU_RF_GROUP3_OID                    {CAPACITY_ROOT,4,10}

#define N2RU_RF_GROUP1_CH1_CH4_OID          {CAPACITY_ROOT,4,11}
#define N2RU_RF_GROUP1_CH5_CH8_OID          {CAPACITY_ROOT,4,12}
#define N2RU_RF_GROUP2_CH1_CH4_OID          {CAPACITY_ROOT,4,13}
#define N2RU_RF_GROUP2_CH5_CH8_OID          {CAPACITY_ROOT,4,14}
#define N2RU_RF_GROUP3_CH1_CH4_OID          {CAPACITY_ROOT,4,15}
#define N2RU_RF_GROUP3_CH5_CH8_OID          {CAPACITY_ROOT,4,16}

#define CAPACITY_GROUP_OID                  {CAPACITY_ROOT,5,1}
#define CAPACITY_RU_CH1_OID                 {CAPACITY_ROOT,5,2}
#define CAPACITY_RU_CH2_OID                 {CAPACITY_ROOT,5,3}
#define CAPACITY_RU_CH3_OID                 {CAPACITY_ROOT,5,4}
#define CAPACITY_RU_CH4_OID                 {CAPACITY_ROOT,5,5}
#define CAPACITY_UPDATE_OID                 {CAPACITY_ROOT,5,6}

/* N2RU */
#define CAPACITY_N2RU_CH1_FREQ_RANGE        {CAPACITY_ROOT,5,7}
#define CAPACITY_N2RU_CH1_OID               {CAPACITY_ROOT,5,8}
#define CAPACITY_N2RU_CH2_FREQ_RANGE        {CAPACITY_ROOT,5,9}
#define CAPACITY_N2RU_CH2_OID               {CAPACITY_ROOT,5,10}
#define CAPACITY_N2RU_CH3_FREQ_RANGE        {CAPACITY_ROOT,5,11}
#define CAPACITY_N2RU_CH3_OID               {CAPACITY_ROOT,5,12}
#define CAPACITY_N2RU_CH4_FREQ_RANGE        {CAPACITY_ROOT,5,13}
#define CAPACITY_N2RU_CH4_OID               {CAPACITY_ROOT,5,14}
#define CAPACITY_N2RU_CH5_FREQ_RANGE        {CAPACITY_ROOT,5,15}
#define CAPACITY_N2RU_CH5_OID               {CAPACITY_ROOT,5,16}
#define CAPACITY_N2RU_CH6_FREQ_RANGE        {CAPACITY_ROOT,5,17}
#define CAPACITY_N2RU_CH6_OID               {CAPACITY_ROOT,5,18}
#define CAPACITY_N2RU_CH7_FREQ_RANGE        {CAPACITY_ROOT,5,19}
#define CAPACITY_N2RU_CH7_OID               {CAPACITY_ROOT,5,20}
#define CAPACITY_N2RU_CH8_FREQ_RANGE        {CAPACITY_ROOT,5,21}
#define CAPACITY_N2RU_CH8_OID               {CAPACITY_ROOT,5,22}

#define N2RU_CAPACITY_UPDATE_OID            {CAPACITY_ROOT,5,23}

#define SERVICE_SUN_TIME_START_OID          {CAPACITY_ROOT,6,1}
#define SERVICE_SUN_TIME_END_OID            {CAPACITY_ROOT,6,2}
#define SERVICE_SUN_WORK_GROUP_OID          {CAPACITY_ROOT,6,3}
#define SERVICE_SUN_NONWORK_GROUP_OID       {CAPACITY_ROOT,6,4}

#define SERVICE_MON_TIME_START_OID          {CAPACITY_ROOT,7,1}
#define SERVICE_MON_TIME_END_OID            {CAPACITY_ROOT,7,2}
#define SERVICE_MON_WORK_GROUP_OID          {CAPACITY_ROOT,7,3}
#define SERVICE_MON_NONWORK_GROUP_OID       {CAPACITY_ROOT,7,4}

#define SERVICE_TUE_TIME_START_OID          {CAPACITY_ROOT,8,1}
#define SERVICE_TUE_TIME_END_OID            {CAPACITY_ROOT,8,2}
#define SERVICE_TUE_WORK_GROUP_OID          {CAPACITY_ROOT,8,3}
#define SERVICE_TUE_NONWORK_GROUP_OID       {CAPACITY_ROOT,8,4}

#define SERVICE_WED_TIME_START_OID          {CAPACITY_ROOT,9,1}
#define SERVICE_WED_TIME_END_OID            {CAPACITY_ROOT,9,2}
#define SERVICE_WED_WORK_GROUP_OID          {CAPACITY_ROOT,9,3}
#define SERVICE_WED_NONWORK_GROUP_OID       {CAPACITY_ROOT,9,4}

#define SERVICE_THU_TIME_START_OID          {CAPACITY_ROOT,10,1}
#define SERVICE_THU_TIME_END_OID            {CAPACITY_ROOT,10,2}
#define SERVICE_THU_WORK_GROUP_OID          {CAPACITY_ROOT,10,3}
#define SERVICE_THU_NONWORK_GROUP_OID       {CAPACITY_ROOT,10,4}

#define SERVICE_FRI_TIME_START_OID          {CAPACITY_ROOT,11,1}
#define SERVICE_FRI_TIME_END_OID            {CAPACITY_ROOT,11,2}
#define SERVICE_FRI_WORK_GROUP_OID          {CAPACITY_ROOT,11,3}
#define SERVICE_FRI_NONWORK_GROUP_OID       {CAPACITY_ROOT,11,4}

#define SERVICE_SAT_TIME_START_OID          {CAPACITY_ROOT,12,1}
#define SERVICE_SAT_TIME_END_OID            {CAPACITY_ROOT,12,2}
#define SERVICE_SAT_WORK_GROUP_OID          {CAPACITY_ROOT,12,3}
#define SERVICE_SAT_NONWORK_GROUP_OID       {CAPACITY_ROOT,12,4}

#define END_OID                             {CAPACITY_ROOT,12,5}


//hpru hp-f-ru

/*last 3 bytes of scalar variables oid*/
/**************告警门限*******************/
#define TEMP_HIGH_THR_OID_HP           {2,2,1}

/*Au*/
#define DN_INPWR_LOW_THR1_OID_HP       {2,2,2}
#define DN_INPWR_LOW_THR2_OID_HP       {2,2,3}
#define DN_INPWR_LOW_THR3_OID_HP       {2,2,4}
#define DN_INPWR_LOW_THR4_OID_HP       {2,2,5}

#define DN_INPWR_HIGH_THR1_OID_HP      {2,2,6}
#define DN_INPWR_HIGH_THR2_OID_HP      {2,2,7}
#define DN_INPWR_HIGH_THR3_OID_HP      {2,2,8}
#define DN_INPWR_HIGH_THR4_OID_HP      {2,2,9}
/*RU*/
#define DN_OUTPWR_LOW_THR1_OID_HP      {2,2,2}
#define DN_OUTPWR_LOW_THR2_OID_HP      {2,2,3}
#define DN_OUTPWR_LOW_THR3_OID_HP      {2,2,4}
#define DN_OUTPWR_LOW_THR4_OID_HP      {2,2,5}

#define DN_OUTPWR_HIGH_THR1_OID_HP     {2,2,6}
#define DN_OUTPWR_HIGH_THR2_OID_HP     {2,2,7}
#define DN_OUTPWR_HIGH_THR3_OID_HP     {2,2,8}
#define DN_OUTPWR_HIGH_THR4_OID_HP     {2,2,9}

/*HRU*/
#define VSWR_THR_OID_HP                {2,2,10}
/*******************************************/

/******************设备基本信息*************/
#define FACTORY_ID_OID_HP              {4,1,1}
#define DEV_MODEL_OID_HP               {4,1,2}
#define DEV_SRRIAL_OID_HP              {4,1,3}
#define ARM_VER_OID_HP                 {4,1,4}
#define SITE_ID_OID_HP                 {4,1,6}
#define DEV_ID_OID_HP                  {4,1,7}
#define DEV_LOCATION_HP                {4,1,8}
#define DEV_SITE_HP                    {4,1,9}
#define DEV_TIME_OID_HP                {4,1,10}
#define SITE_REPORT_OID_HP             {4,1,14}

/*MasterAu*/
#define PROTOCOL_SEL                {4,2,1}
#define NMS_IPADDR                  {4,2,2}
#define NMS_PORT                    {4,2,3} 
#define SECONDARY_NMS_IPADDR        {4,2,4}
#define SECONDARY_NMS_PORT          {4,2,5} 

#define DEVICE_IPADD_OID_HP            {4,2,7}
#define NET_MASK_ADD_OID_HP            {4,2,8}
#define DEFAULT_GWAY_OID_HP            {4,2,9}
#define DEVREV_PORT_HP                 {4,2,10}
#define HEARTTIME_HP                   {4,2,11}    

#define FTP_IP_OID_HP                  {4,3,1}
#define FTP_PORT_OID_HP                {4,3,2}
#define FTP_USER_OID_HP                {4,3,3}
#define FTP_PASSRD_OID_HP              {4,3,4}
#define FTP_DIR_OID_HP                 {4,3,5}
#define FTP_FILENAME_OID_HP            {4,3,6}
#define FTP_TRANS_OID_HP               {4,3,7}

#define SEC_USER_NAME_OID_HP           {4,4,1}
#define AUTH_PROTOCOL_OID_HP           {4,4,2}
#define AUTH_PASSWORD_OID_HP           {4,4,3}
#define PRIV_PROTOCOL_OID_HP           {4,4,4}
#define PRIV_PASSWORD_OID_HP           {4,4,5}
#define USM_EDIT_COMFIRM_OID_HP        {4,4,6}
#define PROTOCOL_SEL_RO_OID_HP         {4,4,7}
#define SNMP_TRAP_PROTOCOL          	{4,4,8}
#define TRAP_IPADD_OID_HP              {4,4,9}
#define TRAP_IPADD2_OID_HP             {4,4,10}
#define TRAP_PORT_OID_HP               {4,4,11}
//#define TRAP_ENGINE1_ID             {4,4,12}
//#define TRAP_ENGINE2_ID             {4,4,13}
#define USM_RESET_OID_HP               {4,4,14}

#define RESEND_ENABLE_OID_HP      {4,4,15}
#define RESEND_INTERVAL_OID_HP   {4,4,16}
#define COMMUNITY_OID_HP              {4,4,17}
#define DELETE_HIS_AlM_OID_HP      {4,4,18}


#define SLV_DEV_NUM_OID_HP             {1,1}
/***********************************************/

/*******************系统信息********************/
/*Au&Ru*/
#define RF_SWITCH1_OID_HP              {5,1,1}
#define UP_GAIN_OFFSET1_OID_HP         {5,1,2}
#define DN_GAIN_OFFSET1_OID_HP         {5,1,3}
#define BAND_WIDE1_OID_HP              {5,1,4} 
#define UP_CENTER_FRE1_OID_HP          {5,1,5}
#define DN_CENTER_FRE1_OID_HP          {5,1,6}
/*Ru*/
#define UP_INPUT1_POWER_OID_HP         {5,1,7}
#define PA_OUTPWR1_OID_HP              {5,1,8}
#define CH1_VSWR_OID_HP                {5,1,9}

#define RF_SWITCH2_OID_HP              {5,2,1}
#define UP_GAIN_OFFSET2_OID_HP         {5,2,2}
#define DN_GAIN_OFFSET2_OID_HP         {5,2,3}
#define BAND_WIDE2_OID_HP              {5,2,4} 
#define UP_CENTER_FRE2_OID_HP          {5,2,5}
#define DN_CENTER_FRE2_OID_HP          {5,2,6}
/*Ru*/
#define UP_INPUT2_POWER_OID_HP         {5,2,7}
#define PA_OUTPWR2_OID_HP              {5,2,8}
#define CH2_VSWR_OID_HP                {5,2,9}

#define RF_SWITCH3_OID_HP              {5,3,1}
#define UP_GAIN_OFFSET3_OID_HP         {5,3,2}
#define DN_GAIN_OFFSET3_OID_HP         {5,3,3}
#define BAND_WIDE3_OID_HP              {5,3,4} 
#define UP_CENTER_FRE3_OID_HP          {5,3,5}
#define DN_CENTER_FRE3_OID_HP          {5,3,6}
/*Ru*/
#define UP_INPUT3_POWER_OID_HP         {5,3,7}
#define PA_OUTPWR3_OID_HP              {5,3,8}

#define RF_SWITCH4_OID_HP              {5,4,1}
#define UP_GAIN_OFFSET4_OID_HP         {5,4,2}
#define DN_GAIN_OFFSET4_OID_HP         {5,4,3}
#define BAND_WIDE4_OID_HP              {5,4,4}
#define UP_CENTER_FRE4_OID_HP          {5,4,5}
#define DN_CENTER_FRE4_OID_HP          {5,4,6}
/*Ru*/
#define UP_INPUT4_POWER_OID_HP         {5,4,7}
#define PA_OUTPWR4_OID_HP              {5,4,8}

#define DEV_TEMP_MAX_OID_HP            {5,5,1}
#define ROUTE_ADD_OID_HP               {5,5,2}

/*AU*/

#define MDL1_BASEBAND_UL_OUTPUT_POWER_HP {5,1,7}
#define MDL2_BASEBAND_UL_OUTPUT_POWER_HP {5,2,7}
#define MDL3_BASEBAND_UL_OUTPUT_POWER_HP {5,3,7}
#define MDL4_BASEBAND_UL_OUTPUT_POWER_HP {5,4,7}

#define MDL1_BASEBAND_DL_INPUT_POWER_HP  {5,1,8}
#define MDL2_BASEBAND_DL_INPUT_POWER_HP  {5,2,8}
#define MDL3_BASEBAND_DL_INPUT_POWER_HP  {5,3,8}
#define MDL4_BASEBAND_DL_INPUT_POWER_HP  {5,4,8}

/*RU*/
#define MDL1_BASEBAND_UL_INPUT_POWER_HP  {5,1,10}
#define MDL2_BASEBAND_UL_INPUT_POWER_HP  {5,2,10}
#define MDL3_BASEBAND_UL_INPUT_POWER_HP  {5,3,10}
#define MDL4_BASEBAND_UL_INPUT_POWER_HP  {5,4,10}

#define MDL1_BASEBAND_DL_OUTPUT_POWER_HP {5,1,11}
#define MDL2_BASEBAND_DL_OUTPUT_POWER_HP {5,2,11}
#define MDL3_BASEBAND_DL_OUTPUT_POWER_HP {5,3,11}
#define MDL4_BASEBAND_DL_OUTPUT_POWER_HP {5,4,11}

/*XP-RU*/
/*#define ALM_FAN_DEV_40W		{5,10,1}
#define FAN1_SPEED_40W			{5,10,2}
#define FAN2_SPEED_40W			{5,10,3}
#define FAN_HIGH_THRESHOLD	{5,10,4}
#define FAN_LOW_THRESHOLD		{5,10,5}*/


#define ARU_PA_SERIAL_NUM_1_HP            {5,1,12}
#define ARU_PA_SERIAL_NUM_2_HP            {5,2,12}
#define ARU_PA_SERIAL_NUM_3_HP            {5,3,12}
#define ARU_PA_SERIAL_NUM_4_HP            {5,4,12}


#define Open_Load_Detect_Switch_HP			{5,5,3}
#define Delay_Compensation_Switch_HP       {7,9,1}
/***********************************************/

/***********************通道配置****************/

/*Master-AU时延*/
#define DELAY_TYPE                  {BAND_CONF_ROOT,1,1}
#define DELAY_ADJ_VAL               {BAND_CONF_ROOT,1,2}
#define DELAY_MEA_VAL               {BAND_CONF_ROOT,1,3}
#define DELAY_CONFIRM               {BAND_CONF_ROOT,1,4}
#define DELAY_CUS_VAL               {BAND_CONF_ROOT,1,5}

/*AU TDD*/
#define TDD_MODE_OID_HP                {BAND_CONF_ROOT,2,1}
#define SYNC_INDICATION_OID_HP         {BAND_CONF_ROOT,2,2}
#define DL_CARR_MAINFREQ_OID_HP        {BAND_CONF_ROOT,2,3}
#define DL_CARR_MINORFREQ_OID_HP       {BAND_CONF_ROOT,2,4}
#define ULDL_CONFIG_MANUAL_OID_HP      {BAND_CONF_ROOT,2,5}
#define SUBFRAMECONF_MANUAL_OID_HP     {BAND_CONF_ROOT,2,6}
#define CP_TYPE_MANUAL_OID_HP          {BAND_CONF_ROOT,2,7}
#define ULDL_CONFIG_AUTO_OID_HP        {BAND_CONF_ROOT,2,8}
#define SUBFRAMECONF_AUTO_OID_HP       {BAND_CONF_ROOT,2,9}
#define CP_TYPE_AUTO_OID_HP            {BAND_CONF_ROOT,2,10}

#define BAND_CONF_MDL_1             (3)
#define BAND1_DATA_VALID_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_1,1}
#define BAND1_ULFREQ_LOW_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_1,2}
#define BAND1_ULFREQ_HIGH_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_1,3}
#define BAND1_DLFREQ_LOW_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_1,4}
#define BAND1_DLFREQ_HIGH_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_1,5}

#define BAND1_SIG1_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_1,6}
#define BAND1_SIG1_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_1,7}
#define BAND1_SIG1_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_1,8}
#define BAND1_SIG1_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_1,9}

#define BAND1_SIG2_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_1,10}
#define BAND1_SIG2_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_1,11}
#define BAND1_SIG2_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_1,12}
#define BAND1_SIG2_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_1,13}

#define BAND1_SIG3_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_1,14}
#define BAND1_SIG3_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_1,15}
#define BAND1_SIG3_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_1,16}
#define BAND1_SIG3_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_1,17}

#define BAND_CONF_MDL_2             (4)
#define BAND2_DATA_VALID_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_2,1}
#define BAND2_ULFREQ_LOW_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_2,2}
#define BAND2_ULFREQ_HIGH_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_2,3}
#define BAND2_DLFREQ_LOW_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_2,4}
#define BAND2_DLFREQ_HIGH_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_2,5}

#define BAND2_SIG1_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_2,6}
#define BAND2_SIG1_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_2,7}
#define BAND2_SIG1_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_2,8}
#define BAND2_SIG1_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_2,9}

#define BAND2_SIG2_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_2,10}
#define BAND2_SIG2_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_2,11}
#define BAND2_SIG2_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_2,12}
#define BAND2_SIG2_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_2,13}

#define BAND2_SIG3_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_2,14}
#define BAND2_SIG3_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_2,15}
#define BAND2_SIG3_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_2,16}
#define BAND2_SIG3_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_2,17}

#define BAND_CONF_MDL_3             (5)
#define BAND3_DATA_VALID_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_3,1}
#define BAND3_ULFREQ_LOW_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_3,2}
#define BAND3_ULFREQ_HIGH_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_3,3}
#define BAND3_DLFREQ_LOW_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_3,4}
#define BAND3_DLFREQ_HIGH_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_3,5}

#define BAND3_SIG1_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_3,6}
#define BAND3_SIG1_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_3,7}
#define BAND3_SIG1_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_3,8}
#define BAND3_SIG1_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_3,9}

#define BAND3_SIG2_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_3,10}
#define BAND3_SIG2_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_3,11}
#define BAND3_SIG2_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_3,12}
#define BAND3_SIG2_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_3,13}

#define BAND3_SIG3_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_3,14}
#define BAND3_SIG3_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_3,15}
#define BAND3_SIG3_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_3,16}
#define BAND3_SIG3_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_3,17}

#define BAND_CONF_MDL_4             (6)
#define BAND4_DATA_VALID_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_4,1}
#define BAND4_ULFREQ_LOW_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_4,2}
#define BAND4_ULFREQ_HIGH_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_4,3}
#define BAND4_DLFREQ_LOW_OID_HP        {BAND_CONF_ROOT,BAND_CONF_MDL_4,4}
#define BAND4_DLFREQ_HIGH_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_4,5}

#define BAND4_SIG1_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_4,6}
#define BAND4_SIG1_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_4,7}
#define BAND4_SIG1_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_4,8}
#define BAND4_SIG1_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_4,9}

#define BAND4_SIG2_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_4,10}
#define BAND4_SIG2_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_4,11}
#define BAND4_SIG2_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_4,12}
#define BAND4_SIG2_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_4,13}

#define BAND4_SIG3_BAND_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_4,14}
#define BAND4_SIG3_TYPE_OID_HP         {BAND_CONF_ROOT,BAND_CONF_MDL_4,15}
#define BAND4_SIG3_ULFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_4,16}
#define BAND4_SIG3_DLFREQ_OID_HP       {BAND_CONF_ROOT,BAND_CONF_MDL_4,17}

/*AU*/
#define BAND_CONFIG_UPDATE_OID_HP      {BAND_CONF_ROOT,7,1}
#define LOC_SIG_BAND                {BAND_CONF_ROOT,7,2}
#define LOC_TRS_BAND                {BAND_CONF_ROOT,7,3}
#define SYS_SIG_BAND                {BAND_CONF_ROOT,7,4}
#define SYS_TRS_BAND                {BAND_CONF_ROOT,7,5}
/*RU*/
#define CHANNEL_MAP1_HP                {BAND_CONF_ROOT,8,1}
#define CHANNEL_MAP2_HP                {BAND_CONF_ROOT,8,2}
#define CHANNEL_MAP3_HP                {BAND_CONF_ROOT,8,3}
#define CHANNEL_MAP4_HP                {BAND_CONF_ROOT,8,4}

#define RU_RF_SWITCH1_OID_HP   {7,3,20}
#define RU_UP_ATT1_OID_HP		{7,3,21}
#define RU_DN_ATT1_OID_HP		{7,3,22}

#define RU_RF_SWITCH2_OID_HP   {7,4,20}
#define RU_UP_ATT2_OID_HP		{7,4,21}
#define RU_DN_ATT2_OID_HP		{7,4,22}

#define RU_RF_SWITCH3_OID_HP   {7,5,20}
#define RU_UP_ATT3_OID_HP		{7,5,21}
#define RU_DN_ATT3_OID_HP		{7,5,22}

#define RU_RF_SWITCH4_OID_HP   {7,6,20}
#define RU_UP_ATT4_OID_HP		{7,6,21}
#define RU_DN_ATT4_OID_HP		{7,6,22}

#if 0
/***********************************************************/
/*****************AU commbiner******************************/
#define POI_CONTROL_MODE                {POI_ROOT,1,1}
#define POI_ADJUST_INTER                {POI_ROOT,1,2}
#define POI_ATT_RESET                   {POI_ROOT,1,3}

#define CHAN1_POI_POWER_SWITCH			{POI_ROOT,2,1}
#define CHAN1PORT1_INPUT_POWER          {POI_ROOT,2,2}
#define CHAN1PORT2_INPUT_POWER          {POI_ROOT,2,3}
#define CHAN1PORT3_INPUT_POWER          {POI_ROOT,2,4}
#define CHAN1PORT4_INPUT_POWER          {POI_ROOT,2,5}
#define CHAN1PORT1_ATTENUATION          {POI_ROOT,2,6}
#define CHAN1PORT2_ATTENUATION          {POI_ROOT,2,7}
#define CHAN1PORT3_ATTENUATION          {POI_ROOT,2,8}
#define CHAN1PORT4_ATTENUATION          {POI_ROOT,2,9}
#define CHAN1PORT1_OFFSET               {POI_ROOT,2,10}
#define CHAN1PORT2_OFFSET               {POI_ROOT,2,11}
#define CHAN1PORT3_OFFSET               {POI_ROOT,2,12}
#define CHAN1PORT4_OFFSET               {POI_ROOT,2,13}
#define CHAN1PORT1_OPERATOR             {POI_ROOT,2,14}
#define CHAN1PORT2_OPERATOR             {POI_ROOT,2,15}
#define CHAN1PORT3_OPERATOR             {POI_ROOT,2,16}
#define CHAN1PORT4_OPERATOR             {POI_ROOT,2,17}
#define CHAN1_POI_SERIAL_NUM		    {POI_ROOT,2,18}

#define CHAN2_POI_POWER_SWITCH			{POI_ROOT,3,1}
#define CHAN2PORT1_INPUT_POWER          {POI_ROOT,3,2}
#define CHAN2PORT2_INPUT_POWER          {POI_ROOT,3,3}
#define CHAN2PORT3_INPUT_POWER          {POI_ROOT,3,4}
#define CHAN2PORT4_INPUT_POWER          {POI_ROOT,3,5}
#define CHAN2PORT1_ATTENUATION          {POI_ROOT,3,6}
#define CHAN2PORT2_ATTENUATION          {POI_ROOT,3,7}
#define CHAN2PORT3_ATTENUATION          {POI_ROOT,3,8}
#define CHAN2PORT4_ATTENUATION          {POI_ROOT,3,9}
#define CHAN2PORT1_OFFSET               {POI_ROOT,3,10}
#define CHAN2PORT2_OFFSET               {POI_ROOT,3,11}
#define CHAN2PORT3_OFFSET               {POI_ROOT,3,12}
#define CHAN2PORT4_OFFSET               {POI_ROOT,3,13}
#define CHAN2PORT1_OPERATOR             {POI_ROOT,3,14}
#define CHAN2PORT2_OPERATOR             {POI_ROOT,3,15}
#define CHAN2PORT3_OPERATOR             {POI_ROOT,3,16}
#define CHAN2PORT4_OPERATOR             {POI_ROOT,3,17}
#define CHAN2_POI_SERIAL_NUM		    {POI_ROOT,3,18}

#define CHAN3_POI_POWER_SWITCH			{POI_ROOT,4,1}
#define CHAN3PORT1_INPUT_POWER          {POI_ROOT,4,2}
#define CHAN3PORT2_INPUT_POWER          {POI_ROOT,4,3}
#define CHAN3PORT3_INPUT_POWER          {POI_ROOT,4,4}
#define CHAN3PORT4_INPUT_POWER          {POI_ROOT,4,5}
#define CHAN3PORT1_ATTENUATION          {POI_ROOT,4,6}
#define CHAN3PORT2_ATTENUATION          {POI_ROOT,4,7}
#define CHAN3PORT3_ATTENUATION          {POI_ROOT,4,8}
#define CHAN3PORT4_ATTENUATION          {POI_ROOT,4,9}
#define CHAN3PORT1_OFFSET               {POI_ROOT,4,10}
#define CHAN3PORT2_OFFSET               {POI_ROOT,4,11}
#define CHAN3PORT3_OFFSET               {POI_ROOT,4,12}
#define CHAN3PORT4_OFFSET               {POI_ROOT,4,13}
#define CHAN3PORT1_OPERATOR             {POI_ROOT,4,14}
#define CHAN3PORT2_OPERATOR             {POI_ROOT,4,15}
#define CHAN3PORT3_OPERATOR             {POI_ROOT,4,16}
#define CHAN3PORT4_OPERATOR             {POI_ROOT,4,17}
#define CHAN3_POI_SERIAL_NUM		    {POI_ROOT,4,18}

#define CHAN4_POI_POWER_SWITCH			{POI_ROOT,5,1}
#define CHAN4PORT1_INPUT_POWER          {POI_ROOT,5,2}
#define CHAN4PORT2_INPUT_POWER          {POI_ROOT,5,3}
#define CHAN4PORT3_INPUT_POWER          {POI_ROOT,5,4}
#define CHAN4PORT4_INPUT_POWER          {POI_ROOT,5,5}
#define CHAN4PORT1_ATTENUATION          {POI_ROOT,5,6}
#define CHAN4PORT2_ATTENUATION          {POI_ROOT,5,7}
#define CHAN4PORT3_ATTENUATION          {POI_ROOT,5,8}
#define CHAN4PORT4_ATTENUATION          {POI_ROOT,5,9}
#define CHAN4PORT1_OFFSET               {POI_ROOT,5,10}
#define CHAN4PORT2_OFFSET               {POI_ROOT,5,11}
#define CHAN4PORT3_OFFSET               {POI_ROOT,5,12}
#define CHAN4PORT4_OFFSET               {POI_ROOT,5,13}
#define CHAN4PORT1_OPERATOR             {POI_ROOT,5,14}
#define CHAN4PORT2_OPERATOR             {POI_ROOT,5,15}
#define CHAN4PORT3_OPERATOR             {POI_ROOT,5,16}
#define CHAN4PORT4_OPERATOR             {POI_ROOT,5,17}
#define CHAN4_POI_SERIAL_NUM		    {POI_ROOT,5,18}
#endif

#define WEBOMT_LOGOUT_TIME		 		{4,1,24}
#define MASTER_AU_WIRELESS_DEBUG_PORT	{4,1,25}
#define SLAVE_LOCAL_DEBUG_PORT	 		{4,1,26}
#define TIME_ZONE	 					{4,1,27}


/********************************************************/

/******************RU module-port-gain*******************/
#define CHANNEL1_PORT1_ULGAIN_HP        {GAIN_ROOT,1,1}
#define CHANNEL1_PORT2_ULGAIN_HP        {GAIN_ROOT,1,2}
#define CHANNEL1_PORT3_ULGAIN_HP        {GAIN_ROOT,1,3}
#define CHANNEL1_PORT4_ULGAIN_HP        {GAIN_ROOT,1,4}
#define CHANNEL1_PORT1_DLGAIN_HP        {GAIN_ROOT,1,5}
#define CHANNEL1_PORT2_DLGAIN_HP        {GAIN_ROOT,1,6}
#define CHANNEL1_PORT3_DLGAIN_HP        {GAIN_ROOT,1,7}
#define CHANNEL1_PORT4_DLGAIN_HP        {GAIN_ROOT,1,8}

#define CHANNEL2_PORT1_ULGAIN_HP        {GAIN_ROOT,2,1}
#define CHANNEL2_PORT2_ULGAIN_HP        {GAIN_ROOT,2,2}
#define CHANNEL2_PORT3_ULGAIN_HP        {GAIN_ROOT,2,3}
#define CHANNEL2_PORT4_ULGAIN_HP        {GAIN_ROOT,2,4}
#define CHANNEL2_PORT1_DLGAIN_HP        {GAIN_ROOT,2,5}
#define CHANNEL2_PORT2_DLGAIN_HP        {GAIN_ROOT,2,6}
#define CHANNEL2_PORT3_DLGAIN_HP        {GAIN_ROOT,2,7}
#define CHANNEL2_PORT4_DLGAIN_HP        {GAIN_ROOT,2,8}

#define CHANNEL3_PORT1_ULGAIN_HP        {GAIN_ROOT,3,1}
#define CHANNEL3_PORT2_ULGAIN_HP        {GAIN_ROOT,3,2}
#define CHANNEL3_PORT3_ULGAIN_HP        {GAIN_ROOT,3,3}
#define CHANNEL3_PORT4_ULGAIN_HP        {GAIN_ROOT,3,4}
#define CHANNEL3_PORT1_DLGAIN_HP        {GAIN_ROOT,3,5}
#define CHANNEL3_PORT2_DLGAIN_HP        {GAIN_ROOT,3,6}
#define CHANNEL3_PORT3_DLGAIN_HP        {GAIN_ROOT,3,7}
#define CHANNEL3_PORT4_DLGAIN_HP        {GAIN_ROOT,3,8}

#define CHANNEL4_PORT1_ULGAIN_HP        {GAIN_ROOT,4,1}
#define CHANNEL4_PORT2_ULGAIN_HP        {GAIN_ROOT,4,2}
#define CHANNEL4_PORT3_ULGAIN_HP        {GAIN_ROOT,4,3}
#define CHANNEL4_PORT4_ULGAIN_HP        {GAIN_ROOT,4,4}
#define CHANNEL4_PORT1_DLGAIN_HP        {GAIN_ROOT,4,5}
#define CHANNEL4_PORT2_DLGAIN_HP        {GAIN_ROOT,4,6}
#define CHANNEL4_PORT3_DLGAIN_HP        {GAIN_ROOT,4,7}
#define CHANNEL4_PORT4_DLGAIN_HP        {GAIN_ROOT,4,8}
/******************************************************/

/*******************RU 容量调度节点***********************/
#define CAPACITY_ALLOCATION_SW_OID_HP          {CAPACITY_ROOT,1,1}
#define AU_1_INFO_UL_OID_HP                    {CAPACITY_ROOT,2,1}
#define AU_1_INFO_DL_OID_HP                    {CAPACITY_ROOT,2,2}
#define AU_2_INFO_UL_OID_HP                    {CAPACITY_ROOT,2,3}
#define AU_2_INFO_DL_OID_HP                    {CAPACITY_ROOT,2,4}
#define AU_3_INFO_UL_OID_HP                    {CAPACITY_ROOT,2,5}
#define AU_3_INFO_DL_OID_HP                    {CAPACITY_ROOT,2,6}
#define AU_4_INFO_UL_OID_HP                    {CAPACITY_ROOT,2,7}
#define AU_4_INFO_DL_OID_HP                    {CAPACITY_ROOT,2,8}
#define SAU1_1_INFO_UL_OID_HP                  {CAPACITY_ROOT,2,9}
#define SAU1_1_INFO_DL_OID_HP                  {CAPACITY_ROOT,2,10}
#define SAU1_2_INFO_UL_OID_HP                  {CAPACITY_ROOT,2,11}
#define SAU1_2_INFO_DL_OID_HP                  {CAPACITY_ROOT,2,12}
#define SAU1_3_INFO_UL_OID_HP                  {CAPACITY_ROOT,2,13}
#define SAU1_3_INFO_DL_OID_HP                  {CAPACITY_ROOT,2,14}
#define SAU1_4_INFO_UL_OID_HP                  {CAPACITY_ROOT,2,15}
#define SAU1_4_INFO_DL_OID_HP                  {CAPACITY_ROOT,2,16}
#define SAU2_1_INFO_UL_OID_HP                  {CAPACITY_ROOT,2,17}
#define SAU2_1_INFO_DL_OID_HP                  {CAPACITY_ROOT,2,18}
#define SAU2_2_INFO_UL_OID_HP                  {CAPACITY_ROOT,2,19}
#define SAU2_2_INFO_DL_OID_HP                  {CAPACITY_ROOT,2,20}
#define SAU2_3_INFO_UL_OID_HP                  {CAPACITY_ROOT,2,21}
#define SAU2_3_INFO_DL_OID_HP                  {CAPACITY_ROOT,2,22}
#define SAU2_4_INFO_UL_OID_HP                  {CAPACITY_ROOT,2,23}
#define SAU2_4_INFO_DL_OID_HP                  {CAPACITY_ROOT,2,24}

#define CH1_INFO_ULDL_OID_HP                   {CAPACITY_ROOT,3,1}
#define CH2_INFO_ULDL_OID_HP                   {CAPACITY_ROOT,3,2}
#define CH3_INFO_ULDL_OID_HP                   {CAPACITY_ROOT,3,3}
#define CH4_INFO_ULDL_OID_HP                   {CAPACITY_ROOT,3,4}

#define HRU_RF_GROUP1_OID_HP                   {CAPACITY_ROOT,4,1}
#define HRU_RF_GROUP2_OID_HP                   {CAPACITY_ROOT,4,2}
#define HRU_RF_GROUP3_OID_HP                   {CAPACITY_ROOT,4,3}

#define MRU_RF_GROUP1_OID_HP                   {CAPACITY_ROOT,4,5}
#define MRU_RF_GROUP2_OID_HP                   {CAPACITY_ROOT,4,6}
#define MRU_RF_GROUP3_OID_HP                   {CAPACITY_ROOT,4,7}

#define RU_RF_GROUP1_OID_HP                    {CAPACITY_ROOT,4,8}
#define RU_RF_GROUP2_OID_HP                    {CAPACITY_ROOT,4,9}
#define RU_RF_GROUP3_OID_HP                    {CAPACITY_ROOT,4,10}


#define CAPACITY_GROUP_OID_HP                  {CAPACITY_ROOT,5,1}
#define CAPACITY_RU_CH1_OID_HP                 {CAPACITY_ROOT,5,2}
#define CAPACITY_RU_CH2_OID_HP                 {CAPACITY_ROOT,5,3}
#define CAPACITY_RU_CH3_OID_HP                 {CAPACITY_ROOT,5,4}
#define CAPACITY_RU_CH4_OID_HP                 {CAPACITY_ROOT,5,5}
#define CAPACITY_UPDATE_OID_HP                 {CAPACITY_ROOT,5,6}

#define SERVICE_SUN_TIME_START_OID_HP          {CAPACITY_ROOT,6,1}
#define SERVICE_SUN_TIME_END_OID_HP            {CAPACITY_ROOT,6,2}
#define SERVICE_SUN_WORK_GROUP_OID_HP          {CAPACITY_ROOT,6,3}
#define SERVICE_SUN_NONWORK_GROUP_OID_HP       {CAPACITY_ROOT,6,4}

#define SERVICE_MON_TIME_START_OID_HP          {CAPACITY_ROOT,7,1}
#define SERVICE_MON_TIME_END_OID_HP            {CAPACITY_ROOT,7,2}
#define SERVICE_MON_WORK_GROUP_OID_HP          {CAPACITY_ROOT,7,3}
#define SERVICE_MON_NONWORK_GROUP_OID_HP       {CAPACITY_ROOT,7,4}

#define SERVICE_TUE_TIME_START_OID_HP          {CAPACITY_ROOT,8,1}
#define SERVICE_TUE_TIME_END_OID_HP            {CAPACITY_ROOT,8,2}
#define SERVICE_TUE_WORK_GROUP_OID_HP          {CAPACITY_ROOT,8,3}
#define SERVICE_TUE_NONWORK_GROUP_OID_HP       {CAPACITY_ROOT,8,4}

#define SERVICE_WED_TIME_START_OID_HP          {CAPACITY_ROOT,9,1}
#define SERVICE_WED_TIME_END_OID_HP            {CAPACITY_ROOT,9,2}
#define SERVICE_WED_WORK_GROUP_OID_HP          {CAPACITY_ROOT,9,3}
#define SERVICE_WED_NONWORK_GROUP_OID_HP       {CAPACITY_ROOT,9,4}

#define SERVICE_THU_TIME_START_OID_HP          {CAPACITY_ROOT,10,1}
#define SERVICE_THU_TIME_END_OID_HP            {CAPACITY_ROOT,10,2}
#define SERVICE_THU_WORK_GROUP_OID_HP          {CAPACITY_ROOT,10,3}
#define SERVICE_THU_NONWORK_GROUP_OID_HP       {CAPACITY_ROOT,10,4}

#define SERVICE_FRI_TIME_START_OID_HP          {CAPACITY_ROOT,11,1}
#define SERVICE_FRI_TIME_END_OID_HP            {CAPACITY_ROOT,11,2}
#define SERVICE_FRI_WORK_GROUP_OID_HP          {CAPACITY_ROOT,11,3}
#define SERVICE_FRI_NONWORK_GROUP_OID_HP       {CAPACITY_ROOT,11,4}

#define SERVICE_SAT_TIME_START_OID_HP          {CAPACITY_ROOT,12,1}
#define SERVICE_SAT_TIME_END_OID_HP            {CAPACITY_ROOT,12,2}
#define SERVICE_SAT_WORK_GROUP_OID_HP          {CAPACITY_ROOT,12,3}
#define SERVICE_SAT_NONWORK_GROUP_OID_HP       {CAPACITY_ROOT,12,4}
/******************************************************/
#define AU_1_UL_INFO_ID_HP 0X0C83
#define AU_2_UL_INFO_ID_HP 0X0C84
#define AU_3_UL_INFO_ID_HP 0X0C85
#define AU_4_UL_INFO_ID_HP 0X0C86

#define SAU1_1_UL_INFO_ID_HP 0X0C93
#define SAU1_2_UL_INFO_ID_HP 0X0C94
#define SAU1_3_UL_INFO_ID_HP 0X0C95
#define SAU1_4_UL_INFO_ID_HP 0X0C96

#define SAU2_1_UL_INFO_ID_HP 0X0C97
#define SAU2_2_UL_INFO_ID_HP 0X0C98
#define SAU2_3_UL_INFO_ID_HP 0X0C99
#define SAU2_4_UL_INFO_ID_HP 0X0C9A

#define AU_1_DL_INFO_ID_HP 0X0D01
#define AU_2_DL_INFO_ID_HP 0X0D02
#define AU_3_DL_INFO_ID_HP 0X0D03
#define AU_4_DL_INFO_ID_HP 0X0D04

#define SAU1_1_DL_INFO_ID_HP 0X0C9B
#define SAU1_2_DL_INFO_ID_HP 0X0C9C
#define SAU1_3_DL_INFO_ID_HP 0X0C9D
#define SAU1_4_DL_INFO_ID_HP 0X0C9E


#define SAU2_1_DL_INFO_ID_HP 0X0C9F
#define SAU2_2_DL_INFO_ID_HP 0X0D10
#define SAU2_3_DL_INFO_ID_HP 0X0D11
#define SAU2_4_DL_INFO_ID_HP 0X0D12



#define SAU1_1_DL_INFO_ID_HP 0X0C9B
#define SAU1_2_DL_INFO_ID_HP 0X0C9C
#define SAU1_3_DL_INFO_ID_HP 0X0C9D
#define SAU1_4_DL_INFO_ID_HP 0X0C9E


#define SAU2_1_DL_INFO_ID_HP 0X0C9F
#define SAU2_2_DL_INFO_ID_HP 0X0D10
#define SAU2_3_DL_INFO_ID_HP 0X0D11
#define SAU2_4_DL_INFO_ID_HP 0X0D12


#define RU_MDL1_INFO_ID_HP 0X0C14
#define RU_MDL2_INFO_ID_HP 0X0C15
#define RU_MDL3_INFO_ID_HP 0X0C1D
#define RU_MDL4_INFO_ID_HP 0X0C10

#define HP_XP_RU_Group1_HP 0X0C11
#define HP_XP_RU_Group2_HP 0X0C12
#define HP_XP_RU_Group3_HP 0X0C13

#define MRU_Group1_HP 0X0C41
#define MRU_Group2_HP 0X0C42
#define MRU_Group3_HP 0X0C43

#define LP_RU_Group1_HP 0X0C90
#define LP_RU_Group2_HP 0X0C91
#define LP_RU_Group3_HP 0X0C92

#define Capacity_Group_Select_ID_HP   0x0c87
#define MDL1_Mapping_Configuration_HP 0X0C18
#define MDL2_Mapping_Configuration_HP 0X0C19
#define MDL3_Mapping_Configuration_HP 0X0C1A
#define MDL4_Mapping_Configuration_HP 0X0C1B




#endif
