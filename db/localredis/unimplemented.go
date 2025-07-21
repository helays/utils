package localredis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// 以下是未实现的redis.UniversalClient接口方法

func (l *LocalCache) Pipeline() redis.Pipeliner {
	return nil
}

func (l *LocalCache) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return nil, nil
}

func (l *LocalCache) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {

	return nil, nil
}

func (l *LocalCache) TxPipeline() redis.Pipeliner {

	return nil
}

func (l *LocalCache) Command(ctx context.Context) *redis.CommandsInfoCmd {

	return nil
}

func (l *LocalCache) CommandList(ctx context.Context, filter *redis.FilterBy) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) CommandGetKeys(ctx context.Context, commands ...interface{}) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) CommandGetKeysAndFlags(ctx context.Context, commands ...interface{}) *redis.KeyFlagsCmd {

	return nil
}

func (l *LocalCache) ClientGetName(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) Echo(ctx context.Context, message interface{}) *redis.StringCmd {

	return nil
}

func (l *LocalCache) Ping(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) Quit(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) Unlink(ctx context.Context, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BgRewriteAOF(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) BgSave(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClientKill(ctx context.Context, ipPort string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClientKillByFilter(ctx context.Context, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ClientList(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) ClientInfo(ctx context.Context) *redis.ClientInfoCmd {

	return nil
}

func (l *LocalCache) ClientPause(ctx context.Context, dur time.Duration) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) ClientUnpause(ctx context.Context) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) ClientID(ctx context.Context) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ClientUnblock(ctx context.Context, id int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ClientUnblockWithError(ctx context.Context, id int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ConfigGet(ctx context.Context, parameter string) *redis.MapStringStringCmd {

	return nil
}

func (l *LocalCache) ConfigResetStat(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ConfigSet(ctx context.Context, parameter, value string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ConfigRewrite(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) DBSize(ctx context.Context) *redis.IntCmd {

	return nil
}

func (l *LocalCache) FlushAll(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FlushAllAsync(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FlushDB(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FlushDBAsync(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) Info(ctx context.Context, section ...string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) LastSave(ctx context.Context) *redis.IntCmd {

	return nil
}

func (l *LocalCache) Save(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) Shutdown(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ShutdownSave(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ShutdownNoSave(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) SlaveOf(ctx context.Context, host, port string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) SlowLogGet(ctx context.Context, num int64) *redis.SlowLogCmd {

	return nil
}

func (l *LocalCache) Time(ctx context.Context) *redis.TimeCmd {

	return nil
}

func (l *LocalCache) DebugObject(ctx context.Context, key string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) MemoryUsage(ctx context.Context, key string, samples ...int) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ModuleLoadex(ctx context.Context, conf *redis.ModuleLoadexConfig) *redis.StringCmd {

	return nil
}

func (l *LocalCache) ACLDryRun(ctx context.Context, username string, command ...interface{}) *redis.StringCmd {

	return nil
}

func (l *LocalCache) ACLLog(ctx context.Context, count int64) *redis.ACLLogCmd {

	return nil
}

func (l *LocalCache) ACLLogReset(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ACLSetUser(ctx context.Context, username string, rules ...string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ACLDelUser(ctx context.Context, username string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ACLList(ctx context.Context) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ACLCat(ctx context.Context) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ACLCatArgs(ctx context.Context, options *redis.ACLCatArgs) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BitPosSpan(ctx context.Context, key string, bit int8, start, end int64, span string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BitField(ctx context.Context, key string, values ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) BitFieldRO(ctx context.Context, key string, values ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) ClusterMyShardID(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) ClusterMyID(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) ClusterSlots(ctx context.Context) *redis.ClusterSlotsCmd {

	return nil
}

func (l *LocalCache) ClusterShards(ctx context.Context) *redis.ClusterShardsCmd {

	return nil
}

func (l *LocalCache) ClusterLinks(ctx context.Context) *redis.ClusterLinksCmd {

	return nil
}

func (l *LocalCache) ClusterNodes(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) ClusterMeet(ctx context.Context, host, port string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterForget(ctx context.Context, nodeID string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterReplicate(ctx context.Context, nodeID string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterResetSoft(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterResetHard(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterInfo(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) ClusterKeySlot(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ClusterCountFailureReports(ctx context.Context, nodeID string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ClusterCountKeysInSlot(ctx context.Context, slot int) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ClusterDelSlots(ctx context.Context, slots ...int) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterDelSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterSaveConfig(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterSlaves(ctx context.Context, nodeID string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ClusterFailover(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterAddSlots(ctx context.Context, slots ...int) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ClusterAddSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ReadOnly(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ReadWrite(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) Dump(ctx context.Context, key string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) Exists(ctx context.Context, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) ExpireTime(ctx context.Context, key string) *redis.DurationCmd {

	return nil
}

func (l *LocalCache) ExpireNX(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) ExpireXX(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) ExpireGT(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) ExpireLT(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) Move(ctx context.Context, key string, db int) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) ObjectFreq(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {

	return nil
}

func (l *LocalCache) Persist(ctx context.Context, key string) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) PExpireTime(ctx context.Context, key string) *redis.DurationCmd {

	return nil
}

func (l *LocalCache) PTTL(ctx context.Context, key string) *redis.DurationCmd {

	return nil
}

func (l *LocalCache) RandomKey(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) RenameNX(ctx context.Context, key, newkey string) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) SortRO(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) SortStore(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {

	return nil
}

func (l *LocalCache) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {

	return nil
}

func (l *LocalCache) Touch(ctx context.Context, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) TTL(ctx context.Context, key string) *redis.DurationCmd {

	return nil
}

func (l *LocalCache) Type(ctx context.Context, key string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) Copy(ctx context.Context, sourceKey string, destKey string, db int, replace bool) *redis.IntCmd {

	return nil
}

func (l *LocalCache) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {

	return nil
}

func (l *LocalCache) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd {

	return nil
}

func (l *LocalCache) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {

	return nil
}

func (l *LocalCache) GeoPos(ctx context.Context, key string, members ...string) *redis.GeoPosCmd {

	return nil
}

func (l *LocalCache) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {

	return nil
}

func (l *LocalCache) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {

	return nil
}

func (l *LocalCache) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {

	return nil
}

func (l *LocalCache) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {

	return nil
}

func (l *LocalCache) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) *redis.GeoSearchLocationCmd {

	return nil
}

func (l *LocalCache) GeoSearchStore(ctx context.Context, key, store string, q *redis.GeoSearchStoreQuery) *redis.IntCmd {

	return nil
}

func (l *LocalCache) GeoDist(ctx context.Context, key string, member1, member2, unit string) *redis.FloatCmd {

	return nil
}

func (l *LocalCache) GeoHash(ctx context.Context, key string, members ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) HExists(ctx context.Context, key, field string) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) HGetDel(ctx context.Context, key string, fields ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) HGetEX(ctx context.Context, key string, fields ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) HGetEXWithArgs(ctx context.Context, key string, options *redis.HGetEXOptions, fields ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {

	return nil
}

func (l *LocalCache) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) HLen(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) HSetEX(ctx context.Context, key string, fieldsAndValues ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) HSetEXWithArgs(ctx context.Context, key string, options *redis.HSetEXOptions, fieldsAndValues ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {

	return nil
}

func (l *LocalCache) HScanNoValues(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {

	return nil
}

func (l *LocalCache) HVals(ctx context.Context, key string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) HRandField(ctx context.Context, key string, count int) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) HRandFieldWithValues(ctx context.Context, key string, count int) *redis.KeyValueSliceCmd {

	return nil
}

func (l *LocalCache) HStrLen(ctx context.Context, key, field string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs redis.HExpireArgs, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HPExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HPExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs redis.HExpireArgs, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs redis.HExpireArgs, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HPExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HPExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs redis.HExpireArgs, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HPersist(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HExpireTime(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HPExpireTime(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HTTL(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) HPTTL(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) PFAdd(ctx context.Context, key string, els ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) PFCount(ctx context.Context, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) PFMerge(ctx context.Context, dest string, keys ...string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) BLMPop(ctx context.Context, timeout time.Duration, direction string, count int64, keys ...string) *redis.KeyValuesCmd {

	return nil
}

func (l *LocalCache) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {

	return nil
}

func (l *LocalCache) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {

	return nil
}

func (l *LocalCache) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) LLen(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) LMPop(ctx context.Context, direction string, count int64, keys ...string) *redis.KeyValuesCmd {

	return nil
}

func (l *LocalCache) LPop(ctx context.Context, key string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) LPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) LPos(ctx context.Context, key string, value string, args redis.LPosArgs) *redis.IntCmd {

	return nil
}

func (l *LocalCache) LPosCount(ctx context.Context, key string, value string, count int64, args redis.LPosArgs) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) LPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) RPop(ctx context.Context, key string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) RPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) RPopLPush(ctx context.Context, source, destination string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) RPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) LMove(ctx context.Context, source, destination, srcpos, destpos string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) *redis.StringCmd {

	return nil
}

func (l *LocalCache) BFAdd(ctx context.Context, key string, element interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) BFCard(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BFExists(ctx context.Context, key string, element interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) BFInfo(ctx context.Context, key string) *redis.BFInfoCmd {

	return nil
}

func (l *LocalCache) BFInfoArg(ctx context.Context, key, option string) *redis.BFInfoCmd {

	return nil
}

func (l *LocalCache) BFInfoCapacity(ctx context.Context, key string) *redis.BFInfoCmd {

	return nil
}

func (l *LocalCache) BFInfoSize(ctx context.Context, key string) *redis.BFInfoCmd {

	return nil
}

func (l *LocalCache) BFInfoFilters(ctx context.Context, key string) *redis.BFInfoCmd {

	return nil
}

func (l *LocalCache) BFInfoItems(ctx context.Context, key string) *redis.BFInfoCmd {

	return nil
}

func (l *LocalCache) BFInfoExpansion(ctx context.Context, key string) *redis.BFInfoCmd {

	return nil
}

func (l *LocalCache) BFInsert(ctx context.Context, key string, options *redis.BFInsertOptions, elements ...interface{}) *redis.BoolSliceCmd {

	return nil
}

func (l *LocalCache) BFMAdd(ctx context.Context, key string, elements ...interface{}) *redis.BoolSliceCmd {

	return nil
}

func (l *LocalCache) BFMExists(ctx context.Context, key string, elements ...interface{}) *redis.BoolSliceCmd {

	return nil
}

func (l *LocalCache) BFReserve(ctx context.Context, key string, errorRate float64, capacity int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) BFReserveExpansion(ctx context.Context, key string, errorRate float64, capacity, expansion int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) BFReserveNonScaling(ctx context.Context, key string, errorRate float64, capacity int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) BFReserveWithArgs(ctx context.Context, key string, options *redis.BFReserveOptions) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) BFScanDump(ctx context.Context, key string, iterator int64) *redis.ScanDumpCmd {

	return nil
}

func (l *LocalCache) BFLoadChunk(ctx context.Context, key string, iterator int64, data interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CFAdd(ctx context.Context, key string, element interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) CFAddNX(ctx context.Context, key string, element interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) CFCount(ctx context.Context, key string, element interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) CFDel(ctx context.Context, key string, element interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) CFExists(ctx context.Context, key string, element interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) CFInfo(ctx context.Context, key string) *redis.CFInfoCmd {

	return nil
}

func (l *LocalCache) CFInsert(ctx context.Context, key string, options *redis.CFInsertOptions, elements ...interface{}) *redis.BoolSliceCmd {

	return nil
}

func (l *LocalCache) CFInsertNX(ctx context.Context, key string, options *redis.CFInsertOptions, elements ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) CFMExists(ctx context.Context, key string, elements ...interface{}) *redis.BoolSliceCmd {

	return nil
}

func (l *LocalCache) CFReserve(ctx context.Context, key string, capacity int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CFReserveWithArgs(ctx context.Context, key string, options *redis.CFReserveOptions) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CFReserveExpansion(ctx context.Context, key string, capacity int64, expansion int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CFReserveBucketSize(ctx context.Context, key string, capacity int64, bucketsize int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CFReserveMaxIterations(ctx context.Context, key string, capacity int64, maxiterations int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CFScanDump(ctx context.Context, key string, iterator int64) *redis.ScanDumpCmd {

	return nil
}

func (l *LocalCache) CFLoadChunk(ctx context.Context, key string, iterator int64, data interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CMSIncrBy(ctx context.Context, key string, elements ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) CMSInfo(ctx context.Context, key string) *redis.CMSInfoCmd {

	return nil
}

func (l *LocalCache) CMSInitByDim(ctx context.Context, key string, width, height int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CMSInitByProb(ctx context.Context, key string, errorRate, probability float64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CMSMerge(ctx context.Context, destKey string, sourceKeys ...string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CMSMergeWithWeight(ctx context.Context, destKey string, sourceKeys map[string]int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) CMSQuery(ctx context.Context, key string, elements ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) TopKAdd(ctx context.Context, key string, elements ...interface{}) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) TopKCount(ctx context.Context, key string, elements ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) TopKIncrBy(ctx context.Context, key string, elements ...interface{}) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) TopKInfo(ctx context.Context, key string) *redis.TopKInfoCmd {

	return nil
}

func (l *LocalCache) TopKList(ctx context.Context, key string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) TopKListWithCount(ctx context.Context, key string) *redis.MapStringIntCmd {

	return nil
}

func (l *LocalCache) TopKQuery(ctx context.Context, key string, elements ...interface{}) *redis.BoolSliceCmd {

	return nil
}

func (l *LocalCache) TopKReserve(ctx context.Context, key string, k int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TopKReserveWithOptions(ctx context.Context, key string, k int64, width, depth int64, decay float64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TDigestAdd(ctx context.Context, key string, elements ...float64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TDigestByRank(ctx context.Context, key string, rank ...uint64) *redis.FloatSliceCmd {

	return nil
}

func (l *LocalCache) TDigestByRevRank(ctx context.Context, key string, rank ...uint64) *redis.FloatSliceCmd {

	return nil
}

func (l *LocalCache) TDigestCDF(ctx context.Context, key string, elements ...float64) *redis.FloatSliceCmd {

	return nil
}

func (l *LocalCache) TDigestCreate(ctx context.Context, key string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TDigestCreateWithCompression(ctx context.Context, key string, compression int64) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TDigestInfo(ctx context.Context, key string) *redis.TDigestInfoCmd {

	return nil
}

func (l *LocalCache) TDigestMax(ctx context.Context, key string) *redis.FloatCmd {

	return nil
}

func (l *LocalCache) TDigestMin(ctx context.Context, key string) *redis.FloatCmd {

	return nil
}

func (l *LocalCache) TDigestMerge(ctx context.Context, destKey string, options *redis.TDigestMergeOptions, sourceKeys ...string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TDigestQuantile(ctx context.Context, key string, elements ...float64) *redis.FloatSliceCmd {

	return nil
}

func (l *LocalCache) TDigestRank(ctx context.Context, key string, values ...float64) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) TDigestReset(ctx context.Context, key string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TDigestRevRank(ctx context.Context, key string, values ...float64) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) TDigestTrimmedMean(ctx context.Context, key string, lowCutQuantile, highCutQuantile float64) *redis.FloatCmd {

	return nil
}

func (l *LocalCache) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) SPublish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) PubSubChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) PubSubNumSub(ctx context.Context, channels ...string) *redis.MapStringIntCmd {

	return nil
}

func (l *LocalCache) PubSubNumPat(ctx context.Context) *redis.IntCmd {

	return nil
}

func (l *LocalCache) PubSubShardChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) PubSubShardNumSub(ctx context.Context, channels ...string) *redis.MapStringIntCmd {

	return nil
}

func (l *LocalCache) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {

	return nil
}

func (l *LocalCache) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {

	return nil
}

func (l *LocalCache) EvalRO(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {

	return nil
}

func (l *LocalCache) EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {

	return nil
}

func (l *LocalCache) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {

	return nil
}

func (l *LocalCache) ScriptFlush(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ScriptKill(ctx context.Context) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FunctionLoad(ctx context.Context, code string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FunctionLoadReplace(ctx context.Context, code string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FunctionDelete(ctx context.Context, libName string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FunctionFlush(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FunctionKill(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FunctionFlushAsync(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FunctionList(ctx context.Context, q redis.FunctionListQuery) *redis.FunctionListCmd {

	return nil
}

func (l *LocalCache) FunctionDump(ctx context.Context) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FunctionRestore(ctx context.Context, libDump string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FunctionStats(ctx context.Context) *redis.FunctionStatsCmd {

	return nil
}

func (l *LocalCache) FCall(ctx context.Context, function string, keys []string, args ...interface{}) *redis.Cmd {

	return nil
}

func (l *LocalCache) FCallRo(ctx context.Context, function string, keys []string, args ...interface{}) *redis.Cmd {

	return nil
}

func (l *LocalCache) FCallRO(ctx context.Context, function string, keys []string, args ...interface{}) *redis.Cmd {

	return nil
}

func (l *LocalCache) FT_List(ctx context.Context) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) FTAggregate(ctx context.Context, index string, query string) *redis.MapStringInterfaceCmd {

	return nil
}

func (l *LocalCache) FTAggregateWithArgs(ctx context.Context, index string, query string, options *redis.FTAggregateOptions) *redis.AggregateCmd {

	return nil
}

func (l *LocalCache) FTAliasAdd(ctx context.Context, index string, alias string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTAliasDel(ctx context.Context, alias string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTAliasUpdate(ctx context.Context, index string, alias string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTAlter(ctx context.Context, index string, skipInitialScan bool, definition []interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTConfigGet(ctx context.Context, option string) *redis.MapMapStringInterfaceCmd {

	return nil
}

func (l *LocalCache) FTConfigSet(ctx context.Context, option string, value interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTCreate(ctx context.Context, index string, options *redis.FTCreateOptions, schema ...*redis.FieldSchema) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTCursorDel(ctx context.Context, index string, cursorId int) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTCursorRead(ctx context.Context, index string, cursorId int, count int) *redis.MapStringInterfaceCmd {

	return nil
}

func (l *LocalCache) FTDictAdd(ctx context.Context, dict string, term ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) FTDictDel(ctx context.Context, dict string, term ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) FTDictDump(ctx context.Context, dict string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) FTDropIndex(ctx context.Context, index string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTDropIndexWithArgs(ctx context.Context, index string, options *redis.FTDropIndexOptions) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTExplain(ctx context.Context, index string, query string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FTExplainWithArgs(ctx context.Context, index string, query string, options *redis.FTExplainOptions) *redis.StringCmd {

	return nil
}

func (l *LocalCache) FTInfo(ctx context.Context, index string) *redis.FTInfoCmd {

	return nil
}

func (l *LocalCache) FTSpellCheck(ctx context.Context, index string, query string) *redis.FTSpellCheckCmd {

	return nil
}

func (l *LocalCache) FTSpellCheckWithArgs(ctx context.Context, index string, query string, options *redis.FTSpellCheckOptions) *redis.FTSpellCheckCmd {

	return nil
}

func (l *LocalCache) FTSearch(ctx context.Context, index string, query string) *redis.FTSearchCmd {

	return nil
}

func (l *LocalCache) FTSearchWithArgs(ctx context.Context, index string, query string, options *redis.FTSearchOptions) *redis.FTSearchCmd {

	return nil
}

func (l *LocalCache) FTSynDump(ctx context.Context, index string) *redis.FTSynDumpCmd {

	return nil
}

func (l *LocalCache) FTSynUpdate(ctx context.Context, index string, synGroupId interface{}, terms []interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTSynUpdateWithArgs(ctx context.Context, index string, synGroupId interface{}, options *redis.FTSynUpdateOptions, terms []interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) FTTagVals(ctx context.Context, index string, field string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) SCard(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) SInterCard(ctx context.Context, limit int64, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) SMIsMember(ctx context.Context, key string, members ...interface{}) *redis.BoolSliceCmd {

	return nil
}

func (l *LocalCache) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {

	return nil
}

func (l *LocalCache) SMove(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) SPop(ctx context.Context, key string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) SRandMember(ctx context.Context, key string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {

	return nil
}

func (l *LocalCache) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {

	return nil
}

func (l *LocalCache) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {

	return nil
}

func (l *LocalCache) BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...string) *redis.ZSliceWithKeyCmd {

	return nil
}

func (l *LocalCache) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZAddLT(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZAddGT(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZAddNX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZAddXX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZAddArgs(ctx context.Context, key string, args redis.ZAddArgs) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZAddArgsIncr(ctx context.Context, key string, args redis.ZAddArgs) *redis.FloatCmd {

	return nil
}

func (l *LocalCache) ZCard(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {

	return nil
}

func (l *LocalCache) ZInter(ctx context.Context, store *redis.ZStore) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZInterWithScores(ctx context.Context, store *redis.ZStore) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZInterCard(ctx context.Context, limit int64, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZInterStore(ctx context.Context, destination string, store *redis.ZStore) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZMPop(ctx context.Context, order string, count int64, keys ...string) *redis.ZSliceWithKeyCmd {

	return nil
}

func (l *LocalCache) ZMScore(ctx context.Context, key string, members ...string) *redis.FloatSliceCmd {

	return nil
}

func (l *LocalCache) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZRangeArgs(ctx context.Context, z redis.ZRangeArgs) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZRangeArgsWithScores(ctx context.Context, z redis.ZRangeArgs) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZRangeStore(ctx context.Context, dst string, z redis.ZRangeArgs) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZRank(ctx context.Context, key, member string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZRankWithScore(ctx context.Context, key, member string) *redis.RankWithScoreCmd {

	return nil
}

func (l *LocalCache) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZRevRankWithScore(ctx context.Context, key, member string) *redis.RankWithScoreCmd {

	return nil
}

func (l *LocalCache) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {

	return nil
}

func (l *LocalCache) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZRandMember(ctx context.Context, key string, count int) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZRandMemberWithScores(ctx context.Context, key string, count int) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZUnion(ctx context.Context, store redis.ZStore) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZUnionWithScores(ctx context.Context, store redis.ZStore) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) ZDiffWithScores(ctx context.Context, keys ...string) *redis.ZSliceCmd {

	return nil
}

func (l *LocalCache) ZDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {

	return nil
}

func (l *LocalCache) Append(ctx context.Context, key, value string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) Decr(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {

	return nil
}

func (l *LocalCache) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {

	return nil
}

func (l *LocalCache) GetEx(ctx context.Context, key string, expiration time.Duration) *redis.StringCmd {

	return nil
}

func (l *LocalCache) GetDel(ctx context.Context, key string) *redis.StringCmd {

	return nil
}

func (l *LocalCache) Incr(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {

	return nil
}

func (l *LocalCache) LCS(ctx context.Context, q *redis.LCSQuery) *redis.LCSCmd {

	return nil
}

func (l *LocalCache) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {

	return nil
}

func (l *LocalCache) MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) SetArgs(ctx context.Context, key string, value interface{}, a redis.SetArgs) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {

	return nil
}

func (l *LocalCache) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) StrLen(ctx context.Context, key string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd {

	return nil
}

func (l *LocalCache) XDel(ctx context.Context, stream string, ids ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XLen(ctx context.Context, stream string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd {

	return nil
}

func (l *LocalCache) XRangeN(ctx context.Context, stream, start, stop string, count int64) *redis.XMessageSliceCmd {

	return nil
}

func (l *LocalCache) XRevRange(ctx context.Context, stream string, start, stop string) *redis.XMessageSliceCmd {

	return nil
}

func (l *LocalCache) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) *redis.XMessageSliceCmd {

	return nil
}

func (l *LocalCache) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {

	return nil
}

func (l *LocalCache) XReadStreams(ctx context.Context, streams ...string) *redis.XStreamSliceCmd {

	return nil
}

func (l *LocalCache) XGroupCreate(ctx context.Context, stream, group, start string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) XGroupSetID(ctx context.Context, stream, group, start string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) XGroupDestroy(ctx context.Context, stream, group string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XGroupCreateConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {

	return nil
}

func (l *LocalCache) XAck(ctx context.Context, stream, group string, ids ...string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XPending(ctx context.Context, stream, group string) *redis.XPendingCmd {

	return nil
}

func (l *LocalCache) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) *redis.XPendingExtCmd {

	return nil
}

func (l *LocalCache) XClaim(ctx context.Context, a *redis.XClaimArgs) *redis.XMessageSliceCmd {

	return nil
}

func (l *LocalCache) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) XAutoClaim(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimCmd {

	return nil
}

func (l *LocalCache) XAutoClaimJustID(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimJustIDCmd {

	return nil
}

func (l *LocalCache) XTrimMaxLen(ctx context.Context, key string, maxLen int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XTrimMinID(ctx context.Context, key string, minID string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) XInfoGroups(ctx context.Context, key string) *redis.XInfoGroupsCmd {

	return nil
}

func (l *LocalCache) XInfoStream(ctx context.Context, key string) *redis.XInfoStreamCmd {

	return nil
}

func (l *LocalCache) XInfoStreamFull(ctx context.Context, key string, count int) *redis.XInfoStreamFullCmd {

	return nil
}

func (l *LocalCache) XInfoConsumers(ctx context.Context, key string, group string) *redis.XInfoConsumersCmd {

	return nil
}

func (l *LocalCache) TSAdd(ctx context.Context, key string, timestamp interface{}, value float64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) TSAddWithArgs(ctx context.Context, key string, timestamp interface{}, value float64, options *redis.TSOptions) *redis.IntCmd {

	return nil
}

func (l *LocalCache) TSCreate(ctx context.Context, key string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TSCreateWithArgs(ctx context.Context, key string, options *redis.TSOptions) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TSAlter(ctx context.Context, key string, options *redis.TSAlterOptions) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TSCreateRule(ctx context.Context, sourceKey string, destKey string, aggregator redis.Aggregator, bucketDuration int) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TSCreateRuleWithArgs(ctx context.Context, sourceKey string, destKey string, aggregator redis.Aggregator, bucketDuration int, options *redis.TSCreateRuleOptions) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TSIncrBy(ctx context.Context, Key string, timestamp float64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) TSIncrByWithArgs(ctx context.Context, key string, timestamp float64, options *redis.TSIncrDecrOptions) *redis.IntCmd {

	return nil
}

func (l *LocalCache) TSDecrBy(ctx context.Context, Key string, timestamp float64) *redis.IntCmd {

	return nil
}

func (l *LocalCache) TSDecrByWithArgs(ctx context.Context, key string, timestamp float64, options *redis.TSIncrDecrOptions) *redis.IntCmd {

	return nil
}

func (l *LocalCache) TSDel(ctx context.Context, Key string, fromTimestamp int, toTimestamp int) *redis.IntCmd {

	return nil
}

func (l *LocalCache) TSDeleteRule(ctx context.Context, sourceKey string, destKey string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) TSGet(ctx context.Context, key string) *redis.TSTimestampValueCmd {

	return nil
}

func (l *LocalCache) TSGetWithArgs(ctx context.Context, key string, options *redis.TSGetOptions) *redis.TSTimestampValueCmd {

	return nil
}

func (l *LocalCache) TSInfo(ctx context.Context, key string) *redis.MapStringInterfaceCmd {

	return nil
}

func (l *LocalCache) TSInfoWithArgs(ctx context.Context, key string, options *redis.TSInfoOptions) *redis.MapStringInterfaceCmd {

	return nil
}

func (l *LocalCache) TSMAdd(ctx context.Context, ktvSlices [][]interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) TSQueryIndex(ctx context.Context, filterExpr []string) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) TSRevRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) *redis.TSTimestampValueSliceCmd {

	return nil
}

func (l *LocalCache) TSRevRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *redis.TSRevRangeOptions) *redis.TSTimestampValueSliceCmd {

	return nil
}

func (l *LocalCache) TSRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) *redis.TSTimestampValueSliceCmd {

	return nil
}

func (l *LocalCache) TSRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *redis.TSRangeOptions) *redis.TSTimestampValueSliceCmd {

	return nil
}

func (l *LocalCache) TSMRange(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string) *redis.MapStringSliceInterfaceCmd {

	return nil
}

func (l *LocalCache) TSMRangeWithArgs(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string, options *redis.TSMRangeOptions) *redis.MapStringSliceInterfaceCmd {

	return nil
}

func (l *LocalCache) TSMRevRange(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string) *redis.MapStringSliceInterfaceCmd {

	return nil
}

func (l *LocalCache) TSMRevRangeWithArgs(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string, options *redis.TSMRevRangeOptions) *redis.MapStringSliceInterfaceCmd {

	return nil
}

func (l *LocalCache) TSMGet(ctx context.Context, filters []string) *redis.MapStringSliceInterfaceCmd {

	return nil
}

func (l *LocalCache) TSMGetWithArgs(ctx context.Context, filters []string, options *redis.TSMGetOptions) *redis.MapStringSliceInterfaceCmd {

	return nil
}

func (l *LocalCache) JSONArrAppend(ctx context.Context, key, path string, values ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) JSONArrIndex(ctx context.Context, key, path string, value ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) JSONArrIndexWithArgs(ctx context.Context, key, path string, options *redis.JSONArrIndexArgs, value ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) JSONArrInsert(ctx context.Context, key, path string, index int64, values ...interface{}) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) JSONArrLen(ctx context.Context, key, path string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) JSONArrPop(ctx context.Context, key, path string, index int) *redis.StringSliceCmd {

	return nil
}

func (l *LocalCache) JSONArrTrim(ctx context.Context, key, path string) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) JSONArrTrimWithArgs(ctx context.Context, key, path string, options *redis.JSONArrTrimArgs) *redis.IntSliceCmd {

	return nil
}

func (l *LocalCache) JSONClear(ctx context.Context, key, path string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) JSONDebugMemory(ctx context.Context, key, path string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) JSONDel(ctx context.Context, key, path string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) JSONForget(ctx context.Context, key, path string) *redis.IntCmd {

	return nil
}

func (l *LocalCache) JSONGet(ctx context.Context, key string, paths ...string) *redis.JSONCmd {

	return nil
}

func (l *LocalCache) JSONGetWithArgs(ctx context.Context, key string, options *redis.JSONGetArgs, paths ...string) *redis.JSONCmd {

	return nil
}

func (l *LocalCache) JSONMerge(ctx context.Context, key, path string, value string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) JSONMSetArgs(ctx context.Context, docs []redis.JSONSetArgs) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) JSONMSet(ctx context.Context, params ...interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) JSONMGet(ctx context.Context, path string, keys ...string) *redis.JSONSliceCmd {

	return nil
}

func (l *LocalCache) JSONNumIncrBy(ctx context.Context, key, path string, value float64) *redis.JSONCmd {

	return nil
}

func (l *LocalCache) JSONObjKeys(ctx context.Context, key, path string) *redis.SliceCmd {

	return nil
}

func (l *LocalCache) JSONObjLen(ctx context.Context, key, path string) *redis.IntPointerSliceCmd {

	return nil
}

func (l *LocalCache) JSONSet(ctx context.Context, key, path string, value interface{}) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) JSONSetMode(ctx context.Context, key, path string, value interface{}, mode string) *redis.StatusCmd {

	return nil
}

func (l *LocalCache) JSONStrAppend(ctx context.Context, key, path, value string) *redis.IntPointerSliceCmd {

	return nil
}

func (l *LocalCache) JSONStrLen(ctx context.Context, key, path string) *redis.IntPointerSliceCmd {

	return nil
}

func (l *LocalCache) JSONToggle(ctx context.Context, key, path string) *redis.IntPointerSliceCmd {
	return nil
}

func (l *LocalCache) JSONType(ctx context.Context, key, path string) *redis.JSONSliceCmd {
	return nil
}

func (l *LocalCache) VAdd(ctx context.Context, key, element string, val redis.Vector) *redis.BoolCmd {
	return nil
}

func (l *LocalCache) VAddWithArgs(ctx context.Context, key, element string, val redis.Vector, addArgs *redis.VAddArgs) *redis.BoolCmd {
	return nil
}

func (l *LocalCache) VCard(ctx context.Context, key string) *redis.IntCmd {
	return nil
}

func (l *LocalCache) VDim(ctx context.Context, key string) *redis.IntCmd {
	return nil
}

func (l *LocalCache) VEmb(ctx context.Context, key, element string, raw bool) *redis.SliceCmd {
	return nil
}

func (l *LocalCache) VGetAttr(ctx context.Context, key, element string) *redis.StringCmd {
	return nil
}

func (l *LocalCache) VInfo(ctx context.Context, key string) *redis.MapStringInterfaceCmd {
	return nil
}

func (l *LocalCache) VLinks(ctx context.Context, key, element string) *redis.StringSliceCmd {
	return nil
}

func (l *LocalCache) VLinksWithScores(ctx context.Context, key, element string) *redis.VectorScoreSliceCmd {
	return nil
}

func (l *LocalCache) VRandMember(ctx context.Context, key string) *redis.StringCmd {
	return nil
}

func (l *LocalCache) VRandMemberCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return nil
}

func (l *LocalCache) VRem(ctx context.Context, key, element string) *redis.BoolCmd {
	return nil
}

func (l *LocalCache) VSetAttr(ctx context.Context, key, element string, attr interface{}) *redis.BoolCmd {
	return nil
}

func (l *LocalCache) VClearAttributes(ctx context.Context, key, element string) *redis.BoolCmd {
	return nil
}

func (l *LocalCache) VSim(ctx context.Context, key string, val redis.Vector) *redis.StringSliceCmd {
	return nil
}

func (l *LocalCache) VSimWithScores(ctx context.Context, key string, val redis.Vector) *redis.VectorScoreSliceCmd {
	return nil
}

func (l *LocalCache) VSimWithArgs(ctx context.Context, key string, val redis.Vector, args *redis.VSimArgs) *redis.StringSliceCmd {
	return nil
}

func (l *LocalCache) VSimWithArgsWithScores(ctx context.Context, key string, val redis.Vector, args *redis.VSimArgs) *redis.VectorScoreSliceCmd {
	return nil
}

func (l *LocalCache) AddHook(hook redis.Hook) {

}

func (l *LocalCache) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	return nil
}

func (l *LocalCache) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	return nil
}

func (l *LocalCache) Process(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (l *LocalCache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return nil
}

func (l *LocalCache) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return nil
}

func (l *LocalCache) SSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return nil
}

func (l *LocalCache) PoolStats() *redis.PoolStats {
	return nil
}
