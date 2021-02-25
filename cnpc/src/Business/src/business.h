#ifndef TAV_BUSINESS_H_
#define TAV_BUSINESS_H_

typedef enum _gms_milestones {
    GMS_MSTONE_ITEM         = 0,
    GMS_MSTONE_PROMO        = 1,
    GMS_MSTONE_PMNT         = 2,
    GMS_MSTONE_FUEL         = 3,
    GMS_MSTONE_NON_FUEL     = 4,
    GMS_MSTONE_TILL_PMNT    = 5,
    GMS_MSTONE_CREDIT       = 6,
    GMS_MSTONE_TILL_PROMO   = 7,
    GMS_MSTONE_PUMP_TIME    = 8,
    GMS_MSTONE_MAP_TILL     = 9,
    GMS_MSTONE_MAP_CAR      = 10,
    GMS_MSTONE_MAX
} gms_milestones;

typedef bool(* query_hander)(void*, void*, void*, void*, void*);

void init_milestones(void*, void*);
bool item_handler(void*, void*, void*, void*, void*);
bool promo_handler(void*, void*, void*, void*, void*);
bool pmnt_handler(void*, void*, void*, void*, void*);
bool fuel_handler(void*, void*, void*, void*, void*);
bool non_fuel_handler(void*, void*, void*, void*, void*);
bool till_promo_handler(void*, void*, void*, void*, void*);
bool till_pmnt_handler(void*, void*, void*, void*, void*);
bool credit_handler(void*, void*, void*, void*, void*);
bool pump_handler(void*, void*, void*, void*, void*);

void bind_till_vehicle(void*, void*, void*, void*);

#endif