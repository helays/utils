package localredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// 以下是未实现的redis.UniversalClient接口方法

// Pipeline
// noinspection all
func (l *LocalCache) Pipeline() redis.Pipeliner { return nil }

// Pipelined
// noinspection all
func (l *LocalCache) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return nil, nil
}

// TxPipelined
// noinspection all
func (l *LocalCache) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return nil, nil
}

// TxPipeline
// noinspection all
func (l *LocalCache) TxPipeline() redis.Pipeliner { return nil }

// Command
// noinspection all
func (l *LocalCache) Command(ctx context.Context) *redis.CommandsInfoCmd { return nil }

// CommandList
// noinspection all
func (l *LocalCache) CommandList(ctx context.Context, filter *redis.FilterBy) *redis.StringSliceCmd {
	return nil
}

// CommandGetKeys
// noinspection all
func (l *LocalCache) CommandGetKeys(ctx context.Context, commands ...interface{}) *redis.StringSliceCmd {
	return nil
}

// CommandGetKeysAndFlags
// noinspection all
func (l *LocalCache) CommandGetKeysAndFlags(ctx context.Context, commands ...interface{}) *redis.KeyFlagsCmd {
	return nil
}

// ClientGetName
// noinspection all
func (l *LocalCache) ClientGetName(ctx context.Context) *redis.StringCmd { return nil }

// Echo
// noinspection all
func (l *LocalCache) Echo(ctx context.Context, message interface{}) *redis.StringCmd { return nil }

// Ping
// noinspection all
func (l *LocalCache) Ping(ctx context.Context) *redis.StatusCmd { return nil }

// Quit
// noinspection all
func (l *LocalCache) Quit(ctx context.Context) *redis.StatusCmd { return nil }

// Unlink
// noinspection all
func (l *LocalCache) Unlink(ctx context.Context, keys ...string) *redis.IntCmd { return nil }

// BgRewriteAOF
// noinspection all
func (l *LocalCache) BgRewriteAOF(ctx context.Context) *redis.StatusCmd { return nil }

// BgSave
// noinspection all
func (l *LocalCache) BgSave(ctx context.Context) *redis.StatusCmd { return nil }

// ClientKill
// noinspection all
func (l *LocalCache) ClientKill(ctx context.Context, ipPort string) *redis.StatusCmd { return nil }

// ClientKillByFilter
// noinspection all
func (l *LocalCache) ClientKillByFilter(ctx context.Context, keys ...string) *redis.IntCmd {
	return nil
}

// ClientList
// noinspection all
func (l *LocalCache) ClientList(ctx context.Context) *redis.StringCmd { return nil }

// ClientInfo
// noinspection all
func (l *LocalCache) ClientInfo(ctx context.Context) *redis.ClientInfoCmd { return nil }

// ClientPause
// noinspection all
func (l *LocalCache) ClientPause(ctx context.Context, dur time.Duration) *redis.BoolCmd { return nil }

// ClientUnpause
// noinspection all
func (l *LocalCache) ClientUnpause(ctx context.Context) *redis.BoolCmd { return nil }

// ClientID
// noinspection all
func (l *LocalCache) ClientID(ctx context.Context) *redis.IntCmd { return nil }

// ClientUnblock
// noinspection all
func (l *LocalCache) ClientUnblock(ctx context.Context, id int64) *redis.IntCmd { return nil }

// ClientUnblockWithError
// noinspection all
func (l *LocalCache) ClientUnblockWithError(ctx context.Context, id int64) *redis.IntCmd { return nil }

// ConfigGet
// noinspection all
func (l *LocalCache) ConfigGet(ctx context.Context, parameter string) *redis.MapStringStringCmd {
	return nil
}

// ConfigResetStat
// noinspection all
func (l *LocalCache) ConfigResetStat(ctx context.Context) *redis.StatusCmd { return nil }

// ConfigSet
// noinspection all
func (l *LocalCache) ConfigSet(ctx context.Context, parameter, value string) *redis.StatusCmd {
	return nil
}

// ConfigRewrite
// noinspection all
func (l *LocalCache) ConfigRewrite(ctx context.Context) *redis.StatusCmd { return nil }

// DBSize
// noinspection all
func (l *LocalCache) DBSize(ctx context.Context) *redis.IntCmd { return nil }

// FlushAll
// noinspection all
func (l *LocalCache) FlushAll(ctx context.Context) *redis.StatusCmd { return nil }

// FlushAllAsync
// noinspection all
func (l *LocalCache) FlushAllAsync(ctx context.Context) *redis.StatusCmd { return nil }

// FlushDB
// noinspection all
func (l *LocalCache) FlushDB(ctx context.Context) *redis.StatusCmd { return nil }

// FlushDBAsync
// noinspection all
func (l *LocalCache) FlushDBAsync(ctx context.Context) *redis.StatusCmd { return nil }

// Info
// noinspection all
func (l *LocalCache) Info(ctx context.Context, section ...string) *redis.StringCmd { return nil }

// LastSave
// noinspection all
func (l *LocalCache) LastSave(ctx context.Context) *redis.IntCmd { return nil }

// Save
// noinspection all
func (l *LocalCache) Save(ctx context.Context) *redis.StatusCmd { return nil }

// Shutdown
// noinspection all
func (l *LocalCache) Shutdown(ctx context.Context) *redis.StatusCmd { return nil }

// ShutdownSave
// noinspection all
func (l *LocalCache) ShutdownSave(ctx context.Context) *redis.StatusCmd { return nil }

// ShutdownNoSave
// noinspection all
func (l *LocalCache) ShutdownNoSave(ctx context.Context) *redis.StatusCmd { return nil }

// SlaveOf
// noinspection all
func (l *LocalCache) SlaveOf(ctx context.Context, host, port string) *redis.StatusCmd { return nil }

// SlowLogGet
// noinspection all
func (l *LocalCache) SlowLogGet(ctx context.Context, num int64) *redis.SlowLogCmd { return nil }

// Time
// noinspection all
func (l *LocalCache) Time(ctx context.Context) *redis.TimeCmd { return nil }

// DebugObject
// noinspection all
func (l *LocalCache) DebugObject(ctx context.Context, key string) *redis.StringCmd { return nil }

// MemoryUsage
// noinspection all
func (l *LocalCache) MemoryUsage(ctx context.Context, key string, samples ...int) *redis.IntCmd {
	return nil
}

// ModuleLoadex
// noinspection all
func (l *LocalCache) ModuleLoadex(ctx context.Context, conf *redis.ModuleLoadexConfig) *redis.StringCmd {
	return nil
}

// ACLDryRun
// noinspection all
func (l *LocalCache) ACLDryRun(ctx context.Context, username string, command ...interface{}) *redis.StringCmd {
	return nil
}

// ACLLog
// noinspection all
func (l *LocalCache) ACLLog(ctx context.Context, count int64) *redis.ACLLogCmd { return nil }

// ACLLogReset
// noinspection all
func (l *LocalCache) ACLLogReset(ctx context.Context) *redis.StatusCmd { return nil }

// ACLSetUser
// noinspection all
func (l *LocalCache) ACLSetUser(ctx context.Context, username string, rules ...string) *redis.StatusCmd {
	return nil
}

// ACLDelUser
// noinspection all
func (l *LocalCache) ACLDelUser(ctx context.Context, username string) *redis.IntCmd { return nil }

// ACLList
// noinspection all
func (l *LocalCache) ACLList(ctx context.Context) *redis.StringSliceCmd { return nil }

// ACLCat
// noinspection all
func (l *LocalCache) ACLCat(ctx context.Context) *redis.StringSliceCmd { return nil }

// ACLCatArgs
// noinspection all
func (l *LocalCache) ACLCatArgs(ctx context.Context, options *redis.ACLCatArgs) *redis.StringSliceCmd {
	return nil
}

// GetBit
// noinspection all
func (l *LocalCache) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd { return nil }

// SetBit
// noinspection all
func (l *LocalCache) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	return nil
}

// BitCount
// noinspection all
func (l *LocalCache) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	return nil
}

// BitOpAnd
// noinspection all
func (l *LocalCache) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return nil
}

// BitOpOr
// noinspection all
func (l *LocalCache) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return nil
}

// BitOpXor
// noinspection all
func (l *LocalCache) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return nil
}

// BitOpNot
// noinspection all
func (l *LocalCache) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	return nil
}

// BitPos
// noinspection all
func (l *LocalCache) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	return nil
}

// BitPosSpan
// noinspection all
func (l *LocalCache) BitPosSpan(ctx context.Context, key string, bit int8, start, end int64, span string) *redis.IntCmd {
	return nil
}

// BitField
// noinspection all
func (l *LocalCache) BitField(ctx context.Context, key string, values ...interface{}) *redis.IntSliceCmd {
	return nil
}

// BitFieldRO
// noinspection all
func (l *LocalCache) BitFieldRO(ctx context.Context, key string, values ...interface{}) *redis.IntSliceCmd {
	return nil
}

// ClusterMyShardID
// noinspection all
func (l *LocalCache) ClusterMyShardID(ctx context.Context) *redis.StringCmd { return nil }

// ClusterMyID
// noinspection all
func (l *LocalCache) ClusterMyID(ctx context.Context) *redis.StringCmd { return nil }

// ClusterSlots
// noinspection all
func (l *LocalCache) ClusterSlots(ctx context.Context) *redis.ClusterSlotsCmd { return nil }

// ClusterShards
// noinspection all
func (l *LocalCache) ClusterShards(ctx context.Context) *redis.ClusterShardsCmd { return nil }

// ClusterLinks
// noinspection all
func (l *LocalCache) ClusterLinks(ctx context.Context) *redis.ClusterLinksCmd { return nil }

// ClusterNodes
// noinspection all
func (l *LocalCache) ClusterNodes(ctx context.Context) *redis.StringCmd { return nil }

// ClusterMeet
// noinspection all
func (l *LocalCache) ClusterMeet(ctx context.Context, host, port string) *redis.StatusCmd { return nil }

// ClusterForget
// noinspection all
func (l *LocalCache) ClusterForget(ctx context.Context, nodeID string) *redis.StatusCmd { return nil }

// ClusterReplicate
// noinspection all
func (l *LocalCache) ClusterReplicate(ctx context.Context, nodeID string) *redis.StatusCmd {
	return nil
}

// ClusterResetSoft
// noinspection all
func (l *LocalCache) ClusterResetSoft(ctx context.Context) *redis.StatusCmd { return nil }

// ClusterResetHard
// noinspection all
func (l *LocalCache) ClusterResetHard(ctx context.Context) *redis.StatusCmd { return nil }

// ClusterInfo
// noinspection all
func (l *LocalCache) ClusterInfo(ctx context.Context) *redis.StringCmd { return nil }

// ClusterKeySlot
// noinspection all
func (l *LocalCache) ClusterKeySlot(ctx context.Context, key string) *redis.IntCmd { return nil }

// ClusterGetKeysInSlot
// noinspection all
func (l *LocalCache) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *redis.StringSliceCmd {
	return nil
}

// ClusterCountFailureReports
// noinspection all
func (l *LocalCache) ClusterCountFailureReports(ctx context.Context, nodeID string) *redis.IntCmd {
	return nil
}

// ClusterCountKeysInSlot
// noinspection all
func (l *LocalCache) ClusterCountKeysInSlot(ctx context.Context, slot int) *redis.IntCmd { return nil }

// ClusterDelSlots
// noinspection all
func (l *LocalCache) ClusterDelSlots(ctx context.Context, slots ...int) *redis.StatusCmd { return nil }

// ClusterDelSlotsRange
// noinspection all
func (l *LocalCache) ClusterDelSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	return nil
}

// ClusterSaveConfig
// noinspection all
func (l *LocalCache) ClusterSaveConfig(ctx context.Context) *redis.StatusCmd { return nil }

// ClusterSlaves
// noinspection all
func (l *LocalCache) ClusterSlaves(ctx context.Context, nodeID string) *redis.StringSliceCmd {
	return nil
}

// ClusterFailover
// noinspection all
func (l *LocalCache) ClusterFailover(ctx context.Context) *redis.StatusCmd { return nil }

// ClusterAddSlots
// noinspection all
func (l *LocalCache) ClusterAddSlots(ctx context.Context, slots ...int) *redis.StatusCmd { return nil }

// ClusterAddSlotsRange
// noinspection all
func (l *LocalCache) ClusterAddSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	return nil
}

// ReadOnly
// noinspection all
func (l *LocalCache) ReadOnly(ctx context.Context) *redis.StatusCmd { return nil }

// ReadWrite
// noinspection all
func (l *LocalCache) ReadWrite(ctx context.Context) *redis.StatusCmd { return nil }

// Dump
// noinspection all
func (l *LocalCache) Dump(ctx context.Context, key string) *redis.StringCmd { return nil }

// Exists
// noinspection all
func (l *LocalCache) Exists(ctx context.Context, keys ...string) *redis.IntCmd { return nil }

// Expire
// noinspection all
func (l *LocalCache) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return nil
}

// ExpireAt
// noinspection all
func (l *LocalCache) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return nil
}

// ExpireTime
// noinspection all
func (l *LocalCache) ExpireTime(ctx context.Context, key string) *redis.DurationCmd { return nil }

// ExpireNX
// noinspection all
func (l *LocalCache) ExpireNX(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return nil
}

// ExpireXX
// noinspection all
func (l *LocalCache) ExpireXX(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return nil
}

// ExpireGT
// noinspection all
func (l *LocalCache) ExpireGT(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return nil
}

// ExpireLT
// noinspection all
func (l *LocalCache) ExpireLT(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return nil
}

// Keys
// noinspection all
func (l *LocalCache) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd { return nil }

// Migrate
// noinspection all
func (l *LocalCache) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	return nil
}

// Move
// noinspection all
func (l *LocalCache) Move(ctx context.Context, key string, db int) *redis.BoolCmd { return nil }

// ObjectFreq
// noinspection all
func (l *LocalCache) ObjectFreq(ctx context.Context, key string) *redis.IntCmd { return nil }

// ObjectRefCount
// noinspection all
func (l *LocalCache) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd { return nil }

// ObjectEncoding
// noinspection all
func (l *LocalCache) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd { return nil }

// ObjectIdleTime
// noinspection all
func (l *LocalCache) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd { return nil }

// Persist
// noinspection all
func (l *LocalCache) Persist(ctx context.Context, key string) *redis.BoolCmd { return nil }

// PExpire
// noinspection all
func (l *LocalCache) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return nil
}

// PExpireAt
// noinspection all
func (l *LocalCache) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return nil
}

// PExpireTime
// noinspection all
func (l *LocalCache) PExpireTime(ctx context.Context, key string) *redis.DurationCmd { return nil }

// PTTL
// noinspection all
func (l *LocalCache) PTTL(ctx context.Context, key string) *redis.DurationCmd { return nil }

// RandomKey
// noinspection all
func (l *LocalCache) RandomKey(ctx context.Context) *redis.StringCmd { return nil }

// Rename
// noinspection all
func (l *LocalCache) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd { return nil }

// RenameNX
// noinspection all
func (l *LocalCache) RenameNX(ctx context.Context, key, newkey string) *redis.BoolCmd { return nil }

// Restore
// noinspection all
func (l *LocalCache) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	return nil
}

// RestoreReplace
// noinspection all
func (l *LocalCache) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	return nil
}

// Sort
// noinspection all
func (l *LocalCache) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	return nil
}

// SortRO
// noinspection all
func (l *LocalCache) SortRO(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	return nil
}

// SortStore
// noinspection all
func (l *LocalCache) SortStore(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {
	return nil
}

// SortInterfaces
// noinspection all
func (l *LocalCache) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	return nil
}

// Touch
// noinspection all
func (l *LocalCache) Touch(ctx context.Context, keys ...string) *redis.IntCmd { return nil }

// TTL
// noinspection all
func (l *LocalCache) TTL(ctx context.Context, key string) *redis.DurationCmd { return nil }

// Type
// noinspection all
func (l *LocalCache) Type(ctx context.Context, key string) *redis.StatusCmd { return nil }

// Copy
// noinspection all
func (l *LocalCache) Copy(ctx context.Context, sourceKey string, destKey string, db int, replace bool) *redis.IntCmd {
	return nil
}

// Scan
// noinspection all
func (l *LocalCache) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return nil
}

// ScanType
// noinspection all
func (l *LocalCache) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd {
	return nil
}

// GeoAdd
// noinspection all
func (l *LocalCache) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	return nil
}

// GeoPos
// noinspection all
func (l *LocalCache) GeoPos(ctx context.Context, key string, members ...string) *redis.GeoPosCmd {
	return nil
}

// GeoRadius
// noinspection all
func (l *LocalCache) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return nil
}

// GeoRadiusStore
// noinspection all
func (l *LocalCache) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {
	return nil
}

// GeoRadiusByMember
// noinspection all
func (l *LocalCache) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return nil
}

// GeoRadiusByMemberStore
// noinspection all
func (l *LocalCache) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {
	return nil
}

// GeoSearch
// noinspection all
func (l *LocalCache) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) *redis.StringSliceCmd {
	return nil
}

// GeoSearchLocation
// noinspection all
func (l *LocalCache) GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) *redis.GeoSearchLocationCmd {
	return nil
}

// GeoSearchStore
// noinspection all
func (l *LocalCache) GeoSearchStore(ctx context.Context, key, store string, q *redis.GeoSearchStoreQuery) *redis.IntCmd {
	return nil
}

// GeoDist
// noinspection all
func (l *LocalCache) GeoDist(ctx context.Context, key string, member1, member2, unit string) *redis.FloatCmd {
	return nil
}

// GeoHash
// noinspection all
func (l *LocalCache) GeoHash(ctx context.Context, key string, members ...string) *redis.StringSliceCmd {
	return nil
}

// HDel
// noinspection all
func (l *LocalCache) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	return nil
}

// HExists
// noinspection all
func (l *LocalCache) HExists(ctx context.Context, key, field string) *redis.BoolCmd { return nil }

// HGetDel
// noinspection all
func (l *LocalCache) HGetDel(ctx context.Context, key string, fields ...string) *redis.StringSliceCmd {
	return nil
}

// HGetEX
// noinspection all
func (l *LocalCache) HGetEX(ctx context.Context, key string, fields ...string) *redis.StringSliceCmd {
	return nil
}

// HGetEXWithArgs
// noinspection all
func (l *LocalCache) HGetEXWithArgs(ctx context.Context, key string, options *redis.HGetEXOptions, fields ...string) *redis.StringSliceCmd {
	return nil
}

// HIncrByFloat
// noinspection all
func (l *LocalCache) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	return nil
}

// HKeys
// noinspection all
func (l *LocalCache) HKeys(ctx context.Context, key string) *redis.StringSliceCmd { return nil }

// HLen
// noinspection all
func (l *LocalCache) HLen(ctx context.Context, key string) *redis.IntCmd { return nil }

// HMSet
// noinspection all
func (l *LocalCache) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return nil
}

// HSetEX
// noinspection all
func (l *LocalCache) HSetEX(ctx context.Context, key string, fieldsAndValues ...string) *redis.IntCmd {
	return nil
}

// HSetEXWithArgs
// noinspection all
func (l *LocalCache) HSetEXWithArgs(ctx context.Context, key string, options *redis.HSetEXOptions, fieldsAndValues ...string) *redis.IntCmd {
	return nil
}

// HSetNX
// noinspection all
func (l *LocalCache) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	return nil
}

// HScan
// noinspection all
func (l *LocalCache) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return nil
}

// HScanNoValues
// noinspection all
func (l *LocalCache) HScanNoValues(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return nil
}

// HVals
// noinspection all
func (l *LocalCache) HVals(ctx context.Context, key string) *redis.StringSliceCmd { return nil }

// HRandField
// noinspection all
func (l *LocalCache) HRandField(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return nil
}

// HRandFieldWithValues
// noinspection all
func (l *LocalCache) HRandFieldWithValues(ctx context.Context, key string, count int) *redis.KeyValueSliceCmd {
	return nil
}

// HStrLen
// noinspection all
func (l *LocalCache) HStrLen(ctx context.Context, key, field string) *redis.IntCmd { return nil }

// HExpire
// noinspection all
func (l *LocalCache) HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HExpireWithArgs
// noinspection all
func (l *LocalCache) HExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs redis.HExpireArgs, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HPExpire
// noinspection all
func (l *LocalCache) HPExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HPExpireWithArgs
// noinspection all
func (l *LocalCache) HPExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs redis.HExpireArgs, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HExpireAt
// noinspection all
func (l *LocalCache) HExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HExpireAtWithArgs
// noinspection all
func (l *LocalCache) HExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs redis.HExpireArgs, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HPExpireAt
// noinspection all
func (l *LocalCache) HPExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HPExpireAtWithArgs
// noinspection all
func (l *LocalCache) HPExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs redis.HExpireArgs, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HPersist
// noinspection all
func (l *LocalCache) HPersist(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HExpireTime
// noinspection all
func (l *LocalCache) HExpireTime(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HPExpireTime
// noinspection all
func (l *LocalCache) HPExpireTime(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HTTL
// noinspection all
func (l *LocalCache) HTTL(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {
	return nil
}

// HPTTL
// noinspection all
func (l *LocalCache) HPTTL(ctx context.Context, key string, fields ...string) *redis.IntSliceCmd {
	return nil
}

// PFAdd
// noinspection all
func (l *LocalCache) PFAdd(ctx context.Context, key string, els ...interface{}) *redis.IntCmd {
	return nil
}

// PFCount
// noinspection all
func (l *LocalCache) PFCount(ctx context.Context, keys ...string) *redis.IntCmd { return nil }

// PFMerge
// noinspection all
func (l *LocalCache) PFMerge(ctx context.Context, dest string, keys ...string) *redis.StatusCmd {
	return nil
}

// BLPop
// noinspection all
func (l *LocalCache) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return nil
}

// BLMPop
// noinspection all
func (l *LocalCache) BLMPop(ctx context.Context, timeout time.Duration, direction string, count int64, keys ...string) *redis.KeyValuesCmd {
	return nil
}

// BRPop
// noinspection all
func (l *LocalCache) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return nil
}

// BRPopLPush
// noinspection all
func (l *LocalCache) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	return nil
}

// LIndex
// noinspection all
func (l *LocalCache) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	return nil
}

// LInsert
// noinspection all
func (l *LocalCache) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	return nil
}

// LInsertBefore
// noinspection all
func (l *LocalCache) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	return nil
}

// LInsertAfter
// noinspection all
func (l *LocalCache) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	return nil
}

// LLen
// noinspection all
func (l *LocalCache) LLen(ctx context.Context, key string) *redis.IntCmd { return nil }

// LMPop
// noinspection all
func (l *LocalCache) LMPop(ctx context.Context, direction string, count int64, keys ...string) *redis.KeyValuesCmd {
	return nil
}

// LPop
// noinspection all
func (l *LocalCache) LPop(ctx context.Context, key string) *redis.StringCmd { return nil }

// LPopCount
// noinspection all
func (l *LocalCache) LPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return nil
}

// LPos
// noinspection all
func (l *LocalCache) LPos(ctx context.Context, key string, value string, args redis.LPosArgs) *redis.IntCmd {
	return nil
}

// LPosCount
// noinspection all
func (l *LocalCache) LPosCount(ctx context.Context, key string, value string, count int64, args redis.LPosArgs) *redis.IntSliceCmd {
	return nil
}

// LPush
// noinspection all
func (l *LocalCache) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return nil
}

// LPushX
// noinspection all
func (l *LocalCache) LPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return nil
}

// LRange
// noinspection all
func (l *LocalCache) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return nil
}

// LRem
// noinspection all
func (l *LocalCache) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	return nil
}

// LSet
// noinspection all
func (l *LocalCache) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	return nil
}

// LTrim
// noinspection all
func (l *LocalCache) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	return nil
}

// RPop
// noinspection all
func (l *LocalCache) RPop(ctx context.Context, key string) *redis.StringCmd { return nil }

// RPopCount
// noinspection all
func (l *LocalCache) RPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return nil
}

// RPopLPush
// noinspection all
func (l *LocalCache) RPopLPush(ctx context.Context, source, destination string) *redis.StringCmd {
	return nil
}

// RPush
// noinspection all
func (l *LocalCache) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return nil
}

// RPushX
// noinspection all
func (l *LocalCache) RPushX(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return nil
}

// LMove
// noinspection all
func (l *LocalCache) LMove(ctx context.Context, source, destination, srcpos, destpos string) *redis.StringCmd {
	return nil
}

// BLMove
// noinspection all
func (l *LocalCache) BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) *redis.StringCmd {
	return nil
}

// BFAdd
// noinspection all
func (l *LocalCache) BFAdd(ctx context.Context, key string, element interface{}) *redis.BoolCmd {
	return nil
}

// BFCard
// noinspection all
func (l *LocalCache) BFCard(ctx context.Context, key string) *redis.IntCmd { return nil }

// BFExists
// noinspection all
func (l *LocalCache) BFExists(ctx context.Context, key string, element interface{}) *redis.BoolCmd {
	return nil
}

// BFInfo
// noinspection all
func (l *LocalCache) BFInfo(ctx context.Context, key string) *redis.BFInfoCmd { return nil }

// BFInfoArg
// noinspection all
func (l *LocalCache) BFInfoArg(ctx context.Context, key, option string) *redis.BFInfoCmd { return nil }

// BFInfoCapacity
// noinspection all
func (l *LocalCache) BFInfoCapacity(ctx context.Context, key string) *redis.BFInfoCmd { return nil }

// BFInfoSize
// noinspection all
func (l *LocalCache) BFInfoSize(ctx context.Context, key string) *redis.BFInfoCmd { return nil }

// BFInfoFilters
// noinspection all
func (l *LocalCache) BFInfoFilters(ctx context.Context, key string) *redis.BFInfoCmd { return nil }

// BFInfoItems
// noinspection all
func (l *LocalCache) BFInfoItems(ctx context.Context, key string) *redis.BFInfoCmd { return nil }

// BFInfoExpansion
// noinspection all
func (l *LocalCache) BFInfoExpansion(ctx context.Context, key string) *redis.BFInfoCmd { return nil }

// BFInsert
// noinspection all
func (l *LocalCache) BFInsert(ctx context.Context, key string, options *redis.BFInsertOptions, elements ...interface{}) *redis.BoolSliceCmd {
	return nil
}

// BFMAdd
// noinspection all
func (l *LocalCache) BFMAdd(ctx context.Context, key string, elements ...interface{}) *redis.BoolSliceCmd {
	return nil
}

// BFMExists
// noinspection all
func (l *LocalCache) BFMExists(ctx context.Context, key string, elements ...interface{}) *redis.BoolSliceCmd {
	return nil
}

// BFReserve
// noinspection all
func (l *LocalCache) BFReserve(ctx context.Context, key string, errorRate float64, capacity int64) *redis.StatusCmd {
	return nil
}

// BFReserveExpansion
// noinspection all
func (l *LocalCache) BFReserveExpansion(ctx context.Context, key string, errorRate float64, capacity, expansion int64) *redis.StatusCmd {
	return nil
}

// BFReserveNonScaling
// noinspection all
func (l *LocalCache) BFReserveNonScaling(ctx context.Context, key string, errorRate float64, capacity int64) *redis.StatusCmd {
	return nil
}

// BFReserveWithArgs
// noinspection all
func (l *LocalCache) BFReserveWithArgs(ctx context.Context, key string, options *redis.BFReserveOptions) *redis.StatusCmd {
	return nil
}

// BFScanDump
// noinspection all
func (l *LocalCache) BFScanDump(ctx context.Context, key string, iterator int64) *redis.ScanDumpCmd {
	return nil
}

// BFLoadChunk
// noinspection all
func (l *LocalCache) BFLoadChunk(ctx context.Context, key string, iterator int64, data interface{}) *redis.StatusCmd {
	return nil
}

// CFAdd
// noinspection all
func (l *LocalCache) CFAdd(ctx context.Context, key string, element interface{}) *redis.BoolCmd {
	return nil
}

// CFAddNX
// noinspection all
func (l *LocalCache) CFAddNX(ctx context.Context, key string, element interface{}) *redis.BoolCmd {
	return nil
}

// CFCount
// noinspection all
func (l *LocalCache) CFCount(ctx context.Context, key string, element interface{}) *redis.IntCmd {
	return nil
}

// CFDel
// noinspection all
func (l *LocalCache) CFDel(ctx context.Context, key string, element interface{}) *redis.BoolCmd {
	return nil
}

// CFExists
// noinspection all
func (l *LocalCache) CFExists(ctx context.Context, key string, element interface{}) *redis.BoolCmd {
	return nil
}

// CFInfo
// noinspection all
func (l *LocalCache) CFInfo(ctx context.Context, key string) *redis.CFInfoCmd { return nil }

// CFInsert
// noinspection all
func (l *LocalCache) CFInsert(ctx context.Context, key string, options *redis.CFInsertOptions, elements ...interface{}) *redis.BoolSliceCmd {
	return nil
}

// CFInsertNX
// noinspection all
func (l *LocalCache) CFInsertNX(ctx context.Context, key string, options *redis.CFInsertOptions, elements ...interface{}) *redis.IntSliceCmd {
	return nil
}

// CFMExists
// noinspection all
func (l *LocalCache) CFMExists(ctx context.Context, key string, elements ...interface{}) *redis.BoolSliceCmd {
	return nil
}

// CFReserve
// noinspection all
func (l *LocalCache) CFReserve(ctx context.Context, key string, capacity int64) *redis.StatusCmd {
	return nil
}

// CFReserveWithArgs
// noinspection all
func (l *LocalCache) CFReserveWithArgs(ctx context.Context, key string, options *redis.CFReserveOptions) *redis.StatusCmd {
	return nil
}

// CFReserveExpansion
// noinspection all
func (l *LocalCache) CFReserveExpansion(ctx context.Context, key string, capacity int64, expansion int64) *redis.StatusCmd {
	return nil
}

// CFReserveBucketSize
// noinspection all
func (l *LocalCache) CFReserveBucketSize(ctx context.Context, key string, capacity int64, bucketsize int64) *redis.StatusCmd {
	return nil
}

// CFReserveMaxIterations
// noinspection all
func (l *LocalCache) CFReserveMaxIterations(ctx context.Context, key string, capacity int64, maxiterations int64) *redis.StatusCmd {
	return nil
}

// CFScanDump
// noinspection all
func (l *LocalCache) CFScanDump(ctx context.Context, key string, iterator int64) *redis.ScanDumpCmd {
	return nil
}

// CFLoadChunk
// noinspection all
func (l *LocalCache) CFLoadChunk(ctx context.Context, key string, iterator int64, data interface{}) *redis.StatusCmd {
	return nil
}

// CMSIncrBy
// noinspection all
func (l *LocalCache) CMSIncrBy(ctx context.Context, key string, elements ...interface{}) *redis.IntSliceCmd {
	return nil
}

// CMSInfo
// noinspection all
func (l *LocalCache) CMSInfo(ctx context.Context, key string) *redis.CMSInfoCmd { return nil }

// CMSInitByDim
// noinspection all
func (l *LocalCache) CMSInitByDim(ctx context.Context, key string, width, height int64) *redis.StatusCmd {
	return nil
}

// CMSInitByProb
// noinspection all
func (l *LocalCache) CMSInitByProb(ctx context.Context, key string, errorRate, probability float64) *redis.StatusCmd {
	return nil
}

// CMSMerge
// noinspection all
func (l *LocalCache) CMSMerge(ctx context.Context, destKey string, sourceKeys ...string) *redis.StatusCmd {
	return nil
}

// CMSMergeWithWeight
// noinspection all
func (l *LocalCache) CMSMergeWithWeight(ctx context.Context, destKey string, sourceKeys map[string]int64) *redis.StatusCmd {
	return nil
}

// CMSQuery
// noinspection all
func (l *LocalCache) CMSQuery(ctx context.Context, key string, elements ...interface{}) *redis.IntSliceCmd {
	return nil
}

// TopKAdd
// noinspection all
func (l *LocalCache) TopKAdd(ctx context.Context, key string, elements ...interface{}) *redis.StringSliceCmd {
	return nil
}

// TopKCount
// noinspection all
func (l *LocalCache) TopKCount(ctx context.Context, key string, elements ...interface{}) *redis.IntSliceCmd {
	return nil
}

// TopKIncrBy
// noinspection all
func (l *LocalCache) TopKIncrBy(ctx context.Context, key string, elements ...interface{}) *redis.StringSliceCmd {
	return nil
}

// TopKInfo
// noinspection all
func (l *LocalCache) TopKInfo(ctx context.Context, key string) *redis.TopKInfoCmd { return nil }

// TopKList
// noinspection all
func (l *LocalCache) TopKList(ctx context.Context, key string) *redis.StringSliceCmd { return nil }

// TopKListWithCount
// noinspection all
func (l *LocalCache) TopKListWithCount(ctx context.Context, key string) *redis.MapStringIntCmd {
	return nil
}

// TopKQuery
// noinspection all
func (l *LocalCache) TopKQuery(ctx context.Context, key string, elements ...interface{}) *redis.BoolSliceCmd {
	return nil
}

// TopKReserve
// noinspection all
func (l *LocalCache) TopKReserve(ctx context.Context, key string, k int64) *redis.StatusCmd {
	return nil
}

// TopKReserveWithOptions
// noinspection all
func (l *LocalCache) TopKReserveWithOptions(ctx context.Context, key string, k int64, width, depth int64, decay float64) *redis.StatusCmd {
	return nil
}

// TDigestAdd
// noinspection all
func (l *LocalCache) TDigestAdd(ctx context.Context, key string, elements ...float64) *redis.StatusCmd {
	return nil
}

// TDigestByRank
// noinspection all
func (l *LocalCache) TDigestByRank(ctx context.Context, key string, rank ...uint64) *redis.FloatSliceCmd {
	return nil
}

// TDigestByRevRank
// noinspection all
func (l *LocalCache) TDigestByRevRank(ctx context.Context, key string, rank ...uint64) *redis.FloatSliceCmd {
	return nil
}

// TDigestCDF
// noinspection all
func (l *LocalCache) TDigestCDF(ctx context.Context, key string, elements ...float64) *redis.FloatSliceCmd {
	return nil
}

// TDigestCreate
// noinspection all
func (l *LocalCache) TDigestCreate(ctx context.Context, key string) *redis.StatusCmd { return nil }

// TDigestCreateWithCompression
// noinspection all
func (l *LocalCache) TDigestCreateWithCompression(ctx context.Context, key string, compression int64) *redis.StatusCmd {
	return nil
}

// TDigestInfo
// noinspection all
func (l *LocalCache) TDigestInfo(ctx context.Context, key string) *redis.TDigestInfoCmd { return nil }

// TDigestMax
// noinspection all
func (l *LocalCache) TDigestMax(ctx context.Context, key string) *redis.FloatCmd { return nil }

// TDigestMin
// noinspection all
func (l *LocalCache) TDigestMin(ctx context.Context, key string) *redis.FloatCmd { return nil }

// TDigestMerge
// noinspection all
func (l *LocalCache) TDigestMerge(ctx context.Context, destKey string, options *redis.TDigestMergeOptions, sourceKeys ...string) *redis.StatusCmd {
	return nil
}

// TDigestQuantile
// noinspection all
func (l *LocalCache) TDigestQuantile(ctx context.Context, key string, elements ...float64) *redis.FloatSliceCmd {
	return nil
}

// TDigestRank
// noinspection all
func (l *LocalCache) TDigestRank(ctx context.Context, key string, values ...float64) *redis.IntSliceCmd {
	return nil
}

// TDigestReset
// noinspection all
func (l *LocalCache) TDigestReset(ctx context.Context, key string) *redis.StatusCmd { return nil }

// TDigestRevRank
// noinspection all
func (l *LocalCache) TDigestRevRank(ctx context.Context, key string, values ...float64) *redis.IntSliceCmd {
	return nil
}

// TDigestTrimmedMean
// noinspection all
func (l *LocalCache) TDigestTrimmedMean(ctx context.Context, key string, lowCutQuantile, highCutQuantile float64) *redis.FloatCmd {
	return nil
}

// Publish
// noinspection all
func (l *LocalCache) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	return nil
}

// SPublish
// noinspection all
func (l *LocalCache) SPublish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	return nil
}

// PubSubChannels
// noinspection all
func (l *LocalCache) PubSubChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return nil
}

// PubSubNumSub
// noinspection all
func (l *LocalCache) PubSubNumSub(ctx context.Context, channels ...string) *redis.MapStringIntCmd {
	return nil
}

// PubSubNumPat
// noinspection all
func (l *LocalCache) PubSubNumPat(ctx context.Context) *redis.IntCmd { return nil }

// PubSubShardChannels
// noinspection all
func (l *LocalCache) PubSubShardChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return nil
}

// PubSubShardNumSub
// noinspection all
func (l *LocalCache) PubSubShardNumSub(ctx context.Context, channels ...string) *redis.MapStringIntCmd {
	return nil
}

// Eval
// noinspection all
func (l *LocalCache) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return nil
}

// EvalSha
// noinspection all
func (l *LocalCache) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return nil
}

// EvalRO
// noinspection all
func (l *LocalCache) EvalRO(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return nil
}

// EvalShaRO
// noinspection all
func (l *LocalCache) EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return nil
}

// ScriptExists
// noinspection all
func (l *LocalCache) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	return nil
}

// ScriptFlush
// noinspection all
func (l *LocalCache) ScriptFlush(ctx context.Context) *redis.StatusCmd { return nil }

// ScriptKill
// noinspection all
func (l *LocalCache) ScriptKill(ctx context.Context) *redis.StatusCmd { return nil }

// ScriptLoad
// noinspection all
func (l *LocalCache) ScriptLoad(ctx context.Context, script string) *redis.StringCmd { return nil }

// FunctionLoad
// noinspection all
func (l *LocalCache) FunctionLoad(ctx context.Context, code string) *redis.StringCmd { return nil }

// FunctionLoadReplace
// noinspection all
func (l *LocalCache) FunctionLoadReplace(ctx context.Context, code string) *redis.StringCmd {
	return nil
}

// FunctionDelete
// noinspection all
func (l *LocalCache) FunctionDelete(ctx context.Context, libName string) *redis.StringCmd { return nil }

// FunctionFlush
// noinspection all
func (l *LocalCache) FunctionFlush(ctx context.Context) *redis.StringCmd { return nil }

// FunctionKill
// noinspection all
func (l *LocalCache) FunctionKill(ctx context.Context) *redis.StringCmd { return nil }

// FunctionFlushAsync
// noinspection all
func (l *LocalCache) FunctionFlushAsync(ctx context.Context) *redis.StringCmd { return nil }

// FunctionList
// noinspection all
func (l *LocalCache) FunctionList(ctx context.Context, q redis.FunctionListQuery) *redis.FunctionListCmd {
	return nil
}

// FunctionDump
// noinspection all
func (l *LocalCache) FunctionDump(ctx context.Context) *redis.StringCmd { return nil }

// FunctionRestore
// noinspection all
func (l *LocalCache) FunctionRestore(ctx context.Context, libDump string) *redis.StringCmd {
	return nil
}

// FunctionStats
// noinspection all
func (l *LocalCache) FunctionStats(ctx context.Context) *redis.FunctionStatsCmd { return nil }

// FCall
// noinspection all
func (l *LocalCache) FCall(ctx context.Context, function string, keys []string, args ...interface{}) *redis.Cmd {
	return nil
}

// FCallRo
// noinspection all
func (l *LocalCache) FCallRo(ctx context.Context, function string, keys []string, args ...interface{}) *redis.Cmd {
	return nil
}

// FCallRO
// noinspection all
func (l *LocalCache) FCallRO(ctx context.Context, function string, keys []string, args ...interface{}) *redis.Cmd {
	return nil
}

// FT_List
// noinspection all
func (l *LocalCache) FT_List(ctx context.Context) *redis.StringSliceCmd { return nil }

// FTAggregate
// noinspection all
func (l *LocalCache) FTAggregate(ctx context.Context, index string, query string) *redis.MapStringInterfaceCmd {
	return nil
}

// FTAggregateWithArgs
// noinspection all
func (l *LocalCache) FTAggregateWithArgs(ctx context.Context, index string, query string, options *redis.FTAggregateOptions) *redis.AggregateCmd {
	return nil
}

// FTAliasAdd
// noinspection all
func (l *LocalCache) FTAliasAdd(ctx context.Context, index string, alias string) *redis.StatusCmd {
	return nil
}

// FTAliasDel
// noinspection all
func (l *LocalCache) FTAliasDel(ctx context.Context, alias string) *redis.StatusCmd { return nil }

// FTAliasUpdate
// noinspection all
func (l *LocalCache) FTAliasUpdate(ctx context.Context, index string, alias string) *redis.StatusCmd {
	return nil
}

// FTAlter
// noinspection all
func (l *LocalCache) FTAlter(ctx context.Context, index string, skipInitialScan bool, definition []interface{}) *redis.StatusCmd {
	return nil
}

// FTConfigGet
// noinspection all
func (l *LocalCache) FTConfigGet(ctx context.Context, option string) *redis.MapMapStringInterfaceCmd {
	return nil
}

// FTConfigSet
// noinspection all
func (l *LocalCache) FTConfigSet(ctx context.Context, option string, value interface{}) *redis.StatusCmd {
	return nil
}

// FTCreate
// noinspection all
func (l *LocalCache) FTCreate(ctx context.Context, index string, options *redis.FTCreateOptions, schema ...*redis.FieldSchema) *redis.StatusCmd {
	return nil
}

// FTCursorDel
// noinspection all
func (l *LocalCache) FTCursorDel(ctx context.Context, index string, cursorId int) *redis.StatusCmd {
	return nil
}

// FTCursorRead
// noinspection all
func (l *LocalCache) FTCursorRead(ctx context.Context, index string, cursorId int, count int) *redis.MapStringInterfaceCmd {
	return nil
}

// FTDictAdd
// noinspection all
func (l *LocalCache) FTDictAdd(ctx context.Context, dict string, term ...interface{}) *redis.IntCmd {
	return nil
}

// FTDictDel
// noinspection all
func (l *LocalCache) FTDictDel(ctx context.Context, dict string, term ...interface{}) *redis.IntCmd {
	return nil
}

// FTDictDump
// noinspection all
func (l *LocalCache) FTDictDump(ctx context.Context, dict string) *redis.StringSliceCmd { return nil }

// FTDropIndex
// noinspection all
func (l *LocalCache) FTDropIndex(ctx context.Context, index string) *redis.StatusCmd { return nil }

// FTDropIndexWithArgs
// noinspection all
func (l *LocalCache) FTDropIndexWithArgs(ctx context.Context, index string, options *redis.FTDropIndexOptions) *redis.StatusCmd {
	return nil
}

// FTExplain
// noinspection all
func (l *LocalCache) FTExplain(ctx context.Context, index string, query string) *redis.StringCmd {
	return nil
}

// FTExplainWithArgs
// noinspection all
func (l *LocalCache) FTExplainWithArgs(ctx context.Context, index string, query string, options *redis.FTExplainOptions) *redis.StringCmd {
	return nil
}

// FTInfo
// noinspection all
func (l *LocalCache) FTInfo(ctx context.Context, index string) *redis.FTInfoCmd { return nil }

// FTSpellCheck
// noinspection all
func (l *LocalCache) FTSpellCheck(ctx context.Context, index string, query string) *redis.FTSpellCheckCmd {
	return nil
}

// FTSpellCheckWithArgs
// noinspection all
func (l *LocalCache) FTSpellCheckWithArgs(ctx context.Context, index string, query string, options *redis.FTSpellCheckOptions) *redis.FTSpellCheckCmd {
	return nil
}

// FTSearch
// noinspection all
func (l *LocalCache) FTSearch(ctx context.Context, index string, query string) *redis.FTSearchCmd {
	return nil
}

// FTSearchWithArgs
// noinspection all
func (l *LocalCache) FTSearchWithArgs(ctx context.Context, index string, query string, options *redis.FTSearchOptions) *redis.FTSearchCmd {
	return nil
}

// FTSynDump
// noinspection all
func (l *LocalCache) FTSynDump(ctx context.Context, index string) *redis.FTSynDumpCmd { return nil }

// FTSynUpdate
// noinspection all
func (l *LocalCache) FTSynUpdate(ctx context.Context, index string, synGroupId interface{}, terms []interface{}) *redis.StatusCmd {
	return nil
}

// FTSynUpdateWithArgs
// noinspection all
func (l *LocalCache) FTSynUpdateWithArgs(ctx context.Context, index string, synGroupId interface{}, options *redis.FTSynUpdateOptions, terms []interface{}) *redis.StatusCmd {
	return nil
}

// FTTagVals
// noinspection all
func (l *LocalCache) FTTagVals(ctx context.Context, index string, field string) *redis.StringSliceCmd {
	return nil
}

// SAdd
// noinspection all
func (l *LocalCache) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return nil
}

// SCard
// noinspection all
func (l *LocalCache) SCard(ctx context.Context, key string) *redis.IntCmd { return nil }

// SDiff
// noinspection all
func (l *LocalCache) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd { return nil }

// SDiffStore
// noinspection all
func (l *LocalCache) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return nil
}

// SInter
// noinspection all
func (l *LocalCache) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd { return nil }

// SInterCard
// noinspection all
func (l *LocalCache) SInterCard(ctx context.Context, limit int64, keys ...string) *redis.IntCmd {
	return nil
}

// SInterStore
// noinspection all
func (l *LocalCache) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return nil
}

// SIsMember
// noinspection all
func (l *LocalCache) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	return nil
}

// SMIsMember
// noinspection all
func (l *LocalCache) SMIsMember(ctx context.Context, key string, members ...interface{}) *redis.BoolSliceCmd {
	return nil
}

// SMembers
// noinspection all
func (l *LocalCache) SMembers(ctx context.Context, key string) *redis.StringSliceCmd { return nil }

// SMembersMap
// noinspection all
func (l *LocalCache) SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {
	return nil
}

// SMove
// noinspection all
func (l *LocalCache) SMove(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {
	return nil
}

// SPop
// noinspection all
func (l *LocalCache) SPop(ctx context.Context, key string) *redis.StringCmd { return nil }

// SPopN
// noinspection all
func (l *LocalCache) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return nil
}

// SRandMember
// noinspection all
func (l *LocalCache) SRandMember(ctx context.Context, key string) *redis.StringCmd { return nil }

// SRandMemberN
// noinspection all
func (l *LocalCache) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return nil
}

// SRem
// noinspection all
func (l *LocalCache) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return nil
}

// SScan
// noinspection all
func (l *LocalCache) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return nil
}

// SUnion
// noinspection all
func (l *LocalCache) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd { return nil }

// SUnionStore
// noinspection all
func (l *LocalCache) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return nil
}

// BZPopMax
// noinspection all
func (l *LocalCache) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return nil
}

// BZPopMin
// noinspection all
func (l *LocalCache) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return nil
}

// BZMPop
// noinspection all
func (l *LocalCache) BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...string) *redis.ZSliceWithKeyCmd {
	return nil
}

// ZAdd
// noinspection all
func (l *LocalCache) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return nil
}

// ZAddLT
// noinspection all
func (l *LocalCache) ZAddLT(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return nil
}

// ZAddGT
// noinspection all
func (l *LocalCache) ZAddGT(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return nil
}

// ZAddNX
// noinspection all
func (l *LocalCache) ZAddNX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return nil
}

// ZAddXX
// noinspection all
func (l *LocalCache) ZAddXX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return nil
}

// ZAddArgs
// noinspection all
func (l *LocalCache) ZAddArgs(ctx context.Context, key string, args redis.ZAddArgs) *redis.IntCmd {
	return nil
}

// ZAddArgsIncr
// noinspection all
func (l *LocalCache) ZAddArgsIncr(ctx context.Context, key string, args redis.ZAddArgs) *redis.FloatCmd {
	return nil
}

// ZCard
// noinspection all
func (l *LocalCache) ZCard(ctx context.Context, key string) *redis.IntCmd { return nil }

// ZCount
// noinspection all
func (l *LocalCache) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd { return nil }

// ZLexCount
// noinspection all
func (l *LocalCache) ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd { return nil }

// ZIncrBy
// noinspection all
func (l *LocalCache) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	return nil
}

// ZInter
// noinspection all
func (l *LocalCache) ZInter(ctx context.Context, store *redis.ZStore) *redis.StringSliceCmd {
	return nil
}

// ZInterWithScores
// noinspection all
func (l *LocalCache) ZInterWithScores(ctx context.Context, store *redis.ZStore) *redis.ZSliceCmd {
	return nil
}

// ZInterCard
// noinspection all
func (l *LocalCache) ZInterCard(ctx context.Context, limit int64, keys ...string) *redis.IntCmd {
	return nil
}

// ZInterStore
// noinspection all
func (l *LocalCache) ZInterStore(ctx context.Context, destination string, store *redis.ZStore) *redis.IntCmd {
	return nil
}

// ZMPop
// noinspection all
func (l *LocalCache) ZMPop(ctx context.Context, order string, count int64, keys ...string) *redis.ZSliceWithKeyCmd {
	return nil
}

// ZMScore
// noinspection all
func (l *LocalCache) ZMScore(ctx context.Context, key string, members ...string) *redis.FloatSliceCmd {
	return nil
}

// ZPopMax
// noinspection all
func (l *LocalCache) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return nil
}

// ZPopMin
// noinspection all
func (l *LocalCache) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return nil
}

// ZRange
// noinspection all
func (l *LocalCache) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return nil
}

// ZRangeWithScores
// noinspection all
func (l *LocalCache) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return nil
}

// ZRangeByScore
// noinspection all
func (l *LocalCache) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return nil
}

// ZRangeByLex
// noinspection all
func (l *LocalCache) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return nil
}

// ZRangeByScoreWithScores
// noinspection all
func (l *LocalCache) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return nil
}

// ZRangeArgs
// noinspection all
func (l *LocalCache) ZRangeArgs(ctx context.Context, z redis.ZRangeArgs) *redis.StringSliceCmd {
	return nil
}

// ZRangeArgsWithScores
// noinspection all
func (l *LocalCache) ZRangeArgsWithScores(ctx context.Context, z redis.ZRangeArgs) *redis.ZSliceCmd {
	return nil
}

// ZRangeStore
// noinspection all
func (l *LocalCache) ZRangeStore(ctx context.Context, dst string, z redis.ZRangeArgs) *redis.IntCmd {
	return nil
}

// ZRank
// noinspection all
func (l *LocalCache) ZRank(ctx context.Context, key, member string) *redis.IntCmd { return nil }

// ZRankWithScore
// noinspection all
func (l *LocalCache) ZRankWithScore(ctx context.Context, key, member string) *redis.RankWithScoreCmd {
	return nil
}

// ZRem
// noinspection all
func (l *LocalCache) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return nil
}

// ZRemRangeByRank
// noinspection all
func (l *LocalCache) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	return nil
}

// ZRemRangeByScore
// noinspection all
func (l *LocalCache) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	return nil
}

// ZRemRangeByLex
// noinspection all
func (l *LocalCache) ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {
	return nil
}

// ZRevRange
// noinspection all
func (l *LocalCache) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return nil
}

// ZRevRangeWithScores
// noinspection all
func (l *LocalCache) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return nil
}

// ZRevRangeByScore
// noinspection all
func (l *LocalCache) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return nil
}

// ZRevRangeByLex
// noinspection all
func (l *LocalCache) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return nil
}

// ZRevRangeByScoreWithScores
// noinspection all
func (l *LocalCache) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return nil
}

// ZRevRank
// noinspection all
func (l *LocalCache) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd { return nil }

// ZRevRankWithScore
// noinspection all
func (l *LocalCache) ZRevRankWithScore(ctx context.Context, key, member string) *redis.RankWithScoreCmd {
	return nil
}

// ZScore
// noinspection all
func (l *LocalCache) ZScore(ctx context.Context, key, member string) *redis.FloatCmd { return nil }

// ZUnionStore
// noinspection all
func (l *LocalCache) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {
	return nil
}

// ZRandMember
// noinspection all
func (l *LocalCache) ZRandMember(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return nil
}

// ZRandMemberWithScores
// noinspection all
func (l *LocalCache) ZRandMemberWithScores(ctx context.Context, key string, count int) *redis.ZSliceCmd {
	return nil
}

// ZUnion
// noinspection all
func (l *LocalCache) ZUnion(ctx context.Context, store redis.ZStore) *redis.StringSliceCmd {
	return nil
}

// ZUnionWithScores
// noinspection all
func (l *LocalCache) ZUnionWithScores(ctx context.Context, store redis.ZStore) *redis.ZSliceCmd {
	return nil
}

// ZDiff
// noinspection all
func (l *LocalCache) ZDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd { return nil }

// ZDiffWithScores
// noinspection all
func (l *LocalCache) ZDiffWithScores(ctx context.Context, keys ...string) *redis.ZSliceCmd {
	return nil
}

// ZDiffStore
// noinspection all
func (l *LocalCache) ZDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return nil
}

// ZScan
// noinspection all
func (l *LocalCache) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return nil
}

// Append
// noinspection all
func (l *LocalCache) Append(ctx context.Context, key, value string) *redis.IntCmd { return nil }

// Decr
// noinspection all
func (l *LocalCache) Decr(ctx context.Context, key string) *redis.IntCmd { return nil }

// DecrBy
// noinspection all
func (l *LocalCache) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	return nil
}

// GetRange
// noinspection all
func (l *LocalCache) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	return nil
}

// GetSet
// noinspection all
func (l *LocalCache) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	return nil
}

// GetEx
// noinspection all
func (l *LocalCache) GetEx(ctx context.Context, key string, expiration time.Duration) *redis.StringCmd {
	return nil
}

// GetDel
// noinspection all
func (l *LocalCache) GetDel(ctx context.Context, key string) *redis.StringCmd { return nil }

// Incr
// noinspection all
func (l *LocalCache) Incr(ctx context.Context, key string) *redis.IntCmd { return nil }

// IncrBy
// noinspection all
func (l *LocalCache) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd { return nil }

// IncrByFloat
// noinspection all
func (l *LocalCache) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	return nil
}

// LCS
// noinspection all
func (l *LocalCache) LCS(ctx context.Context, q *redis.LCSQuery) *redis.LCSCmd { return nil }

// MGet
// noinspection all
func (l *LocalCache) MGet(ctx context.Context, keys ...string) *redis.SliceCmd { return nil }

// MSet
// noinspection all
func (l *LocalCache) MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd { return nil }

// MSetNX
// noinspection all
func (l *LocalCache) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd { return nil }

// SetArgs
// noinspection all
func (l *LocalCache) SetArgs(ctx context.Context, key string, value interface{}, a redis.SetArgs) *redis.StatusCmd {
	return nil
}

// SetEx
// noinspection all
func (l *LocalCache) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return nil
}

// SetNX
// noinspection all
func (l *LocalCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return nil
}

// SetXX
// noinspection all
func (l *LocalCache) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return nil
}

// SetRange
// noinspection all
func (l *LocalCache) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	return nil
}

// StrLen
// noinspection all
func (l *LocalCache) StrLen(ctx context.Context, key string) *redis.IntCmd { return nil }

// XAdd
// noinspection all
func (l *LocalCache) XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd { return nil }

// XDel
// noinspection all
func (l *LocalCache) XDel(ctx context.Context, stream string, ids ...string) *redis.IntCmd {
	return nil
}

// XLen
// noinspection all
func (l *LocalCache) XLen(ctx context.Context, stream string) *redis.IntCmd { return nil }

// XRange
// noinspection all
func (l *LocalCache) XRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd {
	return nil
}

// XRangeN
// noinspection all
func (l *LocalCache) XRangeN(ctx context.Context, stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	return nil
}

// XRevRange
// noinspection all
func (l *LocalCache) XRevRange(ctx context.Context, stream string, start, stop string) *redis.XMessageSliceCmd {
	return nil
}

// XRevRangeN
// noinspection all
func (l *LocalCache) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) *redis.XMessageSliceCmd {
	return nil
}

// XRead
// noinspection all
func (l *LocalCache) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {
	return nil
}

// XReadStreams
// noinspection all
func (l *LocalCache) XReadStreams(ctx context.Context, streams ...string) *redis.XStreamSliceCmd {
	return nil
}

// XGroupCreate
// noinspection all
func (l *LocalCache) XGroupCreate(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	return nil
}

// XGroupCreateMkStream
// noinspection all
func (l *LocalCache) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	return nil
}

// XGroupSetID
// noinspection all
func (l *LocalCache) XGroupSetID(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	return nil
}

// XGroupDestroy
// noinspection all
func (l *LocalCache) XGroupDestroy(ctx context.Context, stream, group string) *redis.IntCmd {
	return nil
}

// XGroupCreateConsumer
// noinspection all
func (l *LocalCache) XGroupCreateConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {
	return nil
}

// XGroupDelConsumer
// noinspection all
func (l *LocalCache) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {
	return nil
}

// XReadGroup
// noinspection all
func (l *LocalCache) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	return nil
}

// XAck
// noinspection all
func (l *LocalCache) XAck(ctx context.Context, stream, group string, ids ...string) *redis.IntCmd {
	return nil
}

// XPending
// noinspection all
func (l *LocalCache) XPending(ctx context.Context, stream, group string) *redis.XPendingCmd {
	return nil
}

// XPendingExt
// noinspection all
func (l *LocalCache) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	return nil
}

// XClaim
// noinspection all
func (l *LocalCache) XClaim(ctx context.Context, a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	return nil
}

// XClaimJustID
// noinspection all
func (l *LocalCache) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) *redis.StringSliceCmd {
	return nil
}

// XAutoClaim
// noinspection all
func (l *LocalCache) XAutoClaim(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimCmd {
	return nil
}

// XAutoClaimJustID
// noinspection all
func (l *LocalCache) XAutoClaimJustID(ctx context.Context, a *redis.XAutoClaimArgs) *redis.XAutoClaimJustIDCmd {
	return nil
}

// XTrimMaxLen
// noinspection all
func (l *LocalCache) XTrimMaxLen(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	return nil
}

// XTrimMaxLenApprox
// noinspection all
func (l *LocalCache) XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) *redis.IntCmd {
	return nil
}

// XTrimMinID
// noinspection all
func (l *LocalCache) XTrimMinID(ctx context.Context, key string, minID string) *redis.IntCmd {
	return nil
}

// XTrimMinIDApprox
// noinspection all
func (l *LocalCache) XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) *redis.IntCmd {
	return nil
}

// XInfoGroups
// noinspection all
func (l *LocalCache) XInfoGroups(ctx context.Context, key string) *redis.XInfoGroupsCmd { return nil }

// XInfoStream
// noinspection all
func (l *LocalCache) XInfoStream(ctx context.Context, key string) *redis.XInfoStreamCmd { return nil }

// XInfoStreamFull
// noinspection all
func (l *LocalCache) XInfoStreamFull(ctx context.Context, key string, count int) *redis.XInfoStreamFullCmd {
	return nil
}

// XInfoConsumers
// noinspection all
func (l *LocalCache) XInfoConsumers(ctx context.Context, key string, group string) *redis.XInfoConsumersCmd {
	return nil
}

// TSAdd
// noinspection all
func (l *LocalCache) TSAdd(ctx context.Context, key string, timestamp interface{}, value float64) *redis.IntCmd {
	return nil
}

// TSAddWithArgs
// noinspection all
func (l *LocalCache) TSAddWithArgs(ctx context.Context, key string, timestamp interface{}, value float64, options *redis.TSOptions) *redis.IntCmd {
	return nil
}

// TSCreate
// noinspection all
func (l *LocalCache) TSCreate(ctx context.Context, key string) *redis.StatusCmd { return nil }

// TSCreateWithArgs
// noinspection all
func (l *LocalCache) TSCreateWithArgs(ctx context.Context, key string, options *redis.TSOptions) *redis.StatusCmd {
	return nil
}

// TSAlter
// noinspection all
func (l *LocalCache) TSAlter(ctx context.Context, key string, options *redis.TSAlterOptions) *redis.StatusCmd {
	return nil
}

// TSCreateRule
// noinspection all
func (l *LocalCache) TSCreateRule(ctx context.Context, sourceKey string, destKey string, aggregator redis.Aggregator, bucketDuration int) *redis.StatusCmd {
	return nil
}

// TSCreateRuleWithArgs
// noinspection all
func (l *LocalCache) TSCreateRuleWithArgs(ctx context.Context, sourceKey string, destKey string, aggregator redis.Aggregator, bucketDuration int, options *redis.TSCreateRuleOptions) *redis.StatusCmd {
	return nil
}

// TSIncrBy
// noinspection all
func (l *LocalCache) TSIncrBy(ctx context.Context, Key string, timestamp float64) *redis.IntCmd {
	return nil
}

// TSIncrByWithArgs
// noinspection all
func (l *LocalCache) TSIncrByWithArgs(ctx context.Context, key string, timestamp float64, options *redis.TSIncrDecrOptions) *redis.IntCmd {
	return nil
}

// TSDecrBy
// noinspection all
func (l *LocalCache) TSDecrBy(ctx context.Context, Key string, timestamp float64) *redis.IntCmd {
	return nil
}

// TSDecrByWithArgs
// noinspection all
func (l *LocalCache) TSDecrByWithArgs(ctx context.Context, key string, timestamp float64, options *redis.TSIncrDecrOptions) *redis.IntCmd {
	return nil
}

// TSDel
// noinspection all
func (l *LocalCache) TSDel(ctx context.Context, Key string, fromTimestamp int, toTimestamp int) *redis.IntCmd {
	return nil
}

// TSDeleteRule
// noinspection all
func (l *LocalCache) TSDeleteRule(ctx context.Context, sourceKey string, destKey string) *redis.StatusCmd {
	return nil
}

// TSGet
// noinspection all
func (l *LocalCache) TSGet(ctx context.Context, key string) *redis.TSTimestampValueCmd { return nil }

// TSGetWithArgs
// noinspection all
func (l *LocalCache) TSGetWithArgs(ctx context.Context, key string, options *redis.TSGetOptions) *redis.TSTimestampValueCmd {
	return nil
}

// TSInfo
// noinspection all
func (l *LocalCache) TSInfo(ctx context.Context, key string) *redis.MapStringInterfaceCmd { return nil }

// TSInfoWithArgs
// noinspection all
func (l *LocalCache) TSInfoWithArgs(ctx context.Context, key string, options *redis.TSInfoOptions) *redis.MapStringInterfaceCmd {
	return nil
}

// TSMAdd
// noinspection all
func (l *LocalCache) TSMAdd(ctx context.Context, ktvSlices [][]interface{}) *redis.IntSliceCmd {
	return nil
}

// TSQueryIndex
// noinspection all
func (l *LocalCache) TSQueryIndex(ctx context.Context, filterExpr []string) *redis.StringSliceCmd {
	return nil
}

// TSRevRange
// noinspection all
func (l *LocalCache) TSRevRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) *redis.TSTimestampValueSliceCmd {
	return nil
}

// TSRevRangeWithArgs
// noinspection all
func (l *LocalCache) TSRevRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *redis.TSRevRangeOptions) *redis.TSTimestampValueSliceCmd {
	return nil
}

// TSRange
// noinspection all
func (l *LocalCache) TSRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) *redis.TSTimestampValueSliceCmd {
	return nil
}

// TSRangeWithArgs
// noinspection all
func (l *LocalCache) TSRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *redis.TSRangeOptions) *redis.TSTimestampValueSliceCmd {
	return nil
}

// TSMRange
// noinspection all
func (l *LocalCache) TSMRange(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string) *redis.MapStringSliceInterfaceCmd {
	return nil
}

// TSMRangeWithArgs
// noinspection all
func (l *LocalCache) TSMRangeWithArgs(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string, options *redis.TSMRangeOptions) *redis.MapStringSliceInterfaceCmd {
	return nil
}

// TSMRevRange
// noinspection all
func (l *LocalCache) TSMRevRange(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string) *redis.MapStringSliceInterfaceCmd {
	return nil
}

// TSMRevRangeWithArgs
// noinspection all
func (l *LocalCache) TSMRevRangeWithArgs(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string, options *redis.TSMRevRangeOptions) *redis.MapStringSliceInterfaceCmd {
	return nil
}

// TSMGet
// noinspection all
func (l *LocalCache) TSMGet(ctx context.Context, filters []string) *redis.MapStringSliceInterfaceCmd {
	return nil
}

// TSMGetWithArgs
// noinspection all
func (l *LocalCache) TSMGetWithArgs(ctx context.Context, filters []string, options *redis.TSMGetOptions) *redis.MapStringSliceInterfaceCmd {
	return nil
}

// JSONArrAppend
// noinspection all
func (l *LocalCache) JSONArrAppend(ctx context.Context, key, path string, values ...interface{}) *redis.IntSliceCmd {
	return nil
}

// JSONArrIndex
// noinspection all
func (l *LocalCache) JSONArrIndex(ctx context.Context, key, path string, value ...interface{}) *redis.IntSliceCmd {
	return nil
}

// JSONArrIndexWithArgs
// noinspection all
func (l *LocalCache) JSONArrIndexWithArgs(ctx context.Context, key, path string, options *redis.JSONArrIndexArgs, value ...interface{}) *redis.IntSliceCmd {
	return nil
}

// JSONArrInsert
// noinspection all
func (l *LocalCache) JSONArrInsert(ctx context.Context, key, path string, index int64, values ...interface{}) *redis.IntSliceCmd {
	return nil
}

// JSONArrLen
// noinspection all
func (l *LocalCache) JSONArrLen(ctx context.Context, key, path string) *redis.IntSliceCmd { return nil }

// JSONArrPop
// noinspection all
func (l *LocalCache) JSONArrPop(ctx context.Context, key, path string, index int) *redis.StringSliceCmd {
	return nil
}

// JSONArrTrim
// noinspection all
func (l *LocalCache) JSONArrTrim(ctx context.Context, key, path string) *redis.IntSliceCmd {
	return nil
}

// JSONArrTrimWithArgs
// noinspection all
func (l *LocalCache) JSONArrTrimWithArgs(ctx context.Context, key, path string, options *redis.JSONArrTrimArgs) *redis.IntSliceCmd {
	return nil
}

// JSONClear
// noinspection all
func (l *LocalCache) JSONClear(ctx context.Context, key, path string) *redis.IntCmd { return nil }

// JSONDebugMemory
// noinspection all
func (l *LocalCache) JSONDebugMemory(ctx context.Context, key, path string) *redis.IntCmd { return nil }

// JSONDel
// noinspection all
func (l *LocalCache) JSONDel(ctx context.Context, key, path string) *redis.IntCmd { return nil }

// JSONForget
// noinspection all
func (l *LocalCache) JSONForget(ctx context.Context, key, path string) *redis.IntCmd { return nil }

// JSONGet
// noinspection all
func (l *LocalCache) JSONGet(ctx context.Context, key string, paths ...string) *redis.JSONCmd {
	return nil
}

// JSONGetWithArgs
// noinspection all
func (l *LocalCache) JSONGetWithArgs(ctx context.Context, key string, options *redis.JSONGetArgs, paths ...string) *redis.JSONCmd {
	return nil
}

// JSONMerge
// noinspection all
func (l *LocalCache) JSONMerge(ctx context.Context, key, path string, value string) *redis.StatusCmd {
	return nil
}

// JSONMSetArgs
// noinspection all
func (l *LocalCache) JSONMSetArgs(ctx context.Context, docs []redis.JSONSetArgs) *redis.StatusCmd {
	return nil
}

// JSONMSet
// noinspection all
func (l *LocalCache) JSONMSet(ctx context.Context, params ...interface{}) *redis.StatusCmd {
	return nil
}

// JSONMGet
// noinspection all
func (l *LocalCache) JSONMGet(ctx context.Context, path string, keys ...string) *redis.JSONSliceCmd {
	return nil
}

// JSONNumIncrBy
// noinspection all
func (l *LocalCache) JSONNumIncrBy(ctx context.Context, key, path string, value float64) *redis.JSONCmd {
	return nil
}

// JSONObjKeys
// noinspection all
func (l *LocalCache) JSONObjKeys(ctx context.Context, key, path string) *redis.SliceCmd { return nil }

// JSONObjLen
// noinspection all
func (l *LocalCache) JSONObjLen(ctx context.Context, key, path string) *redis.IntPointerSliceCmd {
	return nil
}

// JSONSet
// noinspection all
func (l *LocalCache) JSONSet(ctx context.Context, key, path string, value interface{}) *redis.StatusCmd {
	return nil
}

// JSONSetMode
// noinspection all
func (l *LocalCache) JSONSetMode(ctx context.Context, key, path string, value interface{}, mode string) *redis.StatusCmd {
	return nil
}

// JSONStrAppend
// noinspection all
func (l *LocalCache) JSONStrAppend(ctx context.Context, key, path, value string) *redis.IntPointerSliceCmd {
	return nil
}

// JSONStrLen
// noinspection all
func (l *LocalCache) JSONStrLen(ctx context.Context, key, path string) *redis.IntPointerSliceCmd {
	return nil
}

// JSONToggle
// noinspection all
func (l *LocalCache) JSONToggle(ctx context.Context, key, path string) *redis.IntPointerSliceCmd {
	return nil
}

// JSONType
// noinspection all
func (l *LocalCache) JSONType(ctx context.Context, key, path string) *redis.JSONSliceCmd { return nil }

// VAdd
// noinspection all
func (l *LocalCache) VAdd(ctx context.Context, key, element string, val redis.Vector) *redis.BoolCmd {
	return nil
}

// VAddWithArgs
// noinspection all
func (l *LocalCache) VAddWithArgs(ctx context.Context, key, element string, val redis.Vector, addArgs *redis.VAddArgs) *redis.BoolCmd {
	return nil
}

// VCard
// noinspection all
func (l *LocalCache) VCard(ctx context.Context, key string) *redis.IntCmd { return nil }

// VDim
// noinspection all
func (l *LocalCache) VDim(ctx context.Context, key string) *redis.IntCmd { return nil }

// VEmb
// noinspection all
func (l *LocalCache) VEmb(ctx context.Context, key, element string, raw bool) *redis.SliceCmd {
	return nil
}

// VGetAttr
// noinspection all
func (l *LocalCache) VGetAttr(ctx context.Context, key, element string) *redis.StringCmd { return nil }

// VInfo
// noinspection all
func (l *LocalCache) VInfo(ctx context.Context, key string) *redis.MapStringInterfaceCmd { return nil }

// VLinks
// noinspection all
func (l *LocalCache) VLinks(ctx context.Context, key, element string) *redis.StringSliceCmd {
	return nil
}

// VLinksWithScores
// noinspection all
func (l *LocalCache) VLinksWithScores(ctx context.Context, key, element string) *redis.VectorScoreSliceCmd {
	return nil
}

// VRandMember
// noinspection all
func (l *LocalCache) VRandMember(ctx context.Context, key string) *redis.StringCmd { return nil }

// VRandMemberCount
// noinspection all
func (l *LocalCache) VRandMemberCount(ctx context.Context, key string, count int) *redis.StringSliceCmd {
	return nil
}

// VRem
// noinspection all
func (l *LocalCache) VRem(ctx context.Context, key, element string) *redis.BoolCmd { return nil }

// VSetAttr
// noinspection all
func (l *LocalCache) VSetAttr(ctx context.Context, key, element string, attr interface{}) *redis.BoolCmd {
	return nil
}

// VClearAttributes
// noinspection all
func (l *LocalCache) VClearAttributes(ctx context.Context, key, element string) *redis.BoolCmd {
	return nil
}

// VSim
// noinspection all
func (l *LocalCache) VSim(ctx context.Context, key string, val redis.Vector) *redis.StringSliceCmd {
	return nil
}

// VSimWithScores
// noinspection all
func (l *LocalCache) VSimWithScores(ctx context.Context, key string, val redis.Vector) *redis.VectorScoreSliceCmd {
	return nil
}

// VSimWithArgs
// noinspection all
func (l *LocalCache) VSimWithArgs(ctx context.Context, key string, val redis.Vector, args *redis.VSimArgs) *redis.StringSliceCmd {
	return nil
}

// VSimWithArgsWithScores
// noinspection all
func (l *LocalCache) VSimWithArgsWithScores(ctx context.Context, key string, val redis.Vector, args *redis.VSimArgs) *redis.VectorScoreSliceCmd {
	return nil
}

// AddHook
// noinspection all
func (l *LocalCache) AddHook(hook redis.Hook) {

}

// Watch
// noinspection all
func (l *LocalCache) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	return nil
}

// Do
// noinspection all
func (l *LocalCache) Do(ctx context.Context, args ...interface{}) *redis.Cmd { return nil }

// Process
// noinspection all
func (l *LocalCache) Process(ctx context.Context, cmd redis.Cmder) error { return nil }

// Subscribe
// noinspection all
func (l *LocalCache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub { return nil }

// PSubscribe
// noinspection all
func (l *LocalCache) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub { return nil }

// SSubscribe
// noinspection all
func (l *LocalCache) SSubscribe(ctx context.Context, channels ...string) *redis.PubSub { return nil }

// PoolStats
// noinspection all
func (l *LocalCache) PoolStats() *redis.PoolStats { return nil }

// BitOpDiff
// noinspection all
func (l *LocalCache) BitOpDiff(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return nil
}

// BitOpDiff1
// noinspection all
func (l *LocalCache) BitOpDiff1(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return nil
}

// BitOpAndOr
// noinspection all
func (l *LocalCache) BitOpAndOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return nil
}

// BitOpOne
// noinspection all
func (l *LocalCache) BitOpOne(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	return nil
}

// XAckDel
// noinspection all
func (l *LocalCache) XAckDel(ctx context.Context, stream string, group string, mode string, ids ...string) *redis.SliceCmd {
	return nil
}

// XDelEx
// noinspection all
func (l *LocalCache) XDelEx(ctx context.Context, stream string, mode string, ids ...string) *redis.SliceCmd {
	return nil
}

// XTrimMaxLenMode
// noinspection all
func (l *LocalCache) XTrimMaxLenMode(ctx context.Context, key string, maxLen int64, mode string) *redis.IntCmd {
	return nil
}

// XTrimMaxLenApproxMode
// noinspection all
func (l *LocalCache) XTrimMaxLenApproxMode(ctx context.Context, key string, maxLen int64, limit int64, mode string) *redis.IntCmd {
	return nil
}

// XTrimMinIDMode
// noinspection all
func (l *LocalCache) XTrimMinIDMode(ctx context.Context, key string, minID string, mode string) *redis.IntCmd {
	return nil
}

// XTrimMinIDApproxMode
// noinspection all
func (l *LocalCache) XTrimMinIDApproxMode(ctx context.Context, key string, minID string, limit int64, mode string) *redis.IntCmd {
	return nil
}

// ClientMaintNotifications
// noinspection all
func (l *LocalCache) ClientMaintNotifications(ctx context.Context, enabled bool, endpointType string) *redis.StatusCmd {
	return nil
}

// SlowLogLen
// noinspection all
func (l *LocalCache) SlowLogLen(ctx context.Context) *redis.IntCmd { return nil }

// SlowLogReset
// noinspection all
func (l *LocalCache) SlowLogReset(ctx context.Context) *redis.StatusCmd { return nil }

// Latency
// noinspection all
func (l *LocalCache) Latency(ctx context.Context) *redis.LatencyCmd { return nil }

// LatencyReset
// noinspection all
func (l *LocalCache) LatencyReset(ctx context.Context, events ...interface{}) *redis.StatusCmd {
	return nil
}

// ACLGenPass
// noinspection all
func (l *LocalCache) ACLGenPass(ctx context.Context, bit int) *redis.StringCmd { return nil }

// ACLUsers
// noinspection all
func (l *LocalCache) ACLUsers(ctx context.Context) *redis.StringSliceCmd { return nil }

// ACLWhoAmI
// noinspection all
func (l *LocalCache) ACLWhoAmI(ctx context.Context) *redis.StringCmd { return nil }

// FTHybrid
// noinspection all
func (l *LocalCache) FTHybrid(ctx context.Context, index string, searchExpr string, vectorField string, vectorData redis.Vector) *redis.FTHybridCmd {
	return nil
}

// FTHybridWithArgs
// noinspection all
func (l *LocalCache) FTHybridWithArgs(ctx context.Context, index string, options *redis.FTHybridOptions) *redis.FTHybridCmd {
	return nil
}

// DelExArgs
// noinspection all
func (l *LocalCache) DelExArgs(ctx context.Context, key string, a redis.DelExArgs) *redis.IntCmd {
	return nil
}

// Digest
// noinspection all
func (l *LocalCache) Digest(ctx context.Context, key string) *redis.DigestCmd { return nil }

// MSetEX
// noinspection all
func (l *LocalCache) MSetEX(ctx context.Context, args redis.MSetEXArgs, values ...interface{}) *redis.IntCmd {
	return nil
}

// SetIFEQ
// noinspection all
func (l *LocalCache) SetIFEQ(ctx context.Context, key string, value interface{}, matchValue interface{}, expiration time.Duration) *redis.StatusCmd {
	return nil
}

// SetIFEQGet
// noinspection all
func (l *LocalCache) SetIFEQGet(ctx context.Context, key string, value interface{}, matchValue interface{}, expiration time.Duration) *redis.StringCmd {
	return nil
}

// SetIFNE
// noinspection all
func (l *LocalCache) SetIFNE(ctx context.Context, key string, value interface{}, matchValue interface{}, expiration time.Duration) *redis.StatusCmd {
	return nil
}

// SetIFNEGet
// noinspection all
func (l *LocalCache) SetIFNEGet(ctx context.Context, key string, value interface{}, matchValue interface{}, expiration time.Duration) *redis.StringCmd {
	return nil
}

// SetIFDEQ
// noinspection all
func (l *LocalCache) SetIFDEQ(ctx context.Context, key string, value interface{}, matchDigest uint64, expiration time.Duration) *redis.StatusCmd {
	return nil
}

// SetIFDEQGet
// noinspection all
func (l *LocalCache) SetIFDEQGet(ctx context.Context, key string, value interface{}, matchDigest uint64, expiration time.Duration) *redis.StringCmd {
	return nil
}

// SetIFDNE
// noinspection all
func (l *LocalCache) SetIFDNE(ctx context.Context, key string, value interface{}, matchDigest uint64, expiration time.Duration) *redis.StatusCmd {
	return nil
}

// SetIFDNEGet
// noinspection all
func (l *LocalCache) SetIFDNEGet(ctx context.Context, key string, value interface{}, matchDigest uint64, expiration time.Duration) *redis.StringCmd {
	return nil
}

// VRange
// noinspection all
func (l *LocalCache) VRange(ctx context.Context, key string, start string, end string, count int64) *redis.StringSliceCmd {
	return nil
}
