package datastruct

import (
	"math"
	"github.com/magicgravity/chatserver/common"
)

const (
	ZSKIPLIST_MAXLEVEL  = 32 	/* Should be enough for 2^32 elements */
	ZSKIPLIST_P = 0.25			/* Skiplist P = 1/4 */
)

type ZSkiplist struct {
	header *ZSkiplistNode
	tail *ZSkiplistNode
	length uint64
	level int
}

type zskiplistLevel struct {
	forward *ZSkiplistNode
	span uint64
}

type ZSkiplistNode struct {
	score float64
	backward *ZSkiplistNode
	levels []zskiplistLevel
	ele string
}

/* Create a skiplist node with the specified number of levels.
 * The SDS string 'ele' is referenced by the node after the call. */
func CreateZSkiplistNode(level int64,score float64,elm string)*ZSkiplistNode{
	znode := ZSkiplistNode{}
	znode.score = score
	znode.levels = make([]zskiplistLevel,level)
	znode.ele = elm
	return &znode
}


/* Create a new skiplist. */
func CreateZSkiplist()*ZSkiplist{
	zsk  := ZSkiplist{}
	zsk.level = 1
	zsk.length = 0
	zsk.header = CreateZSkiplistNode(math.MaxInt64,0,"")

	var j = 0
	for j<ZSKIPLIST_MAXLEVEL {
		zsk.header.levels[j].forward = nil
		zsk.header.levels[j].span = 0
		j++
	}
	zsk.header.backward = nil
	zsk.tail = nil

	return &zsk
}

/* Returns a random level for the new skiplist node we are going to create.
 * The return value of this function is between 1 and ZSKIPLIST_MAXLEVEL
 * (both inclusive), with a powerlaw-alike distribution where higher
 * levels are less likely to be returned. */
func randomLevel()int{
	level := 1
	ranV := common.GenRandomIntWithNegative( -90,32767)
	for (ranV&0xFFFF) < (ZSKIPLIST_P * 0xFFFF) {
		level++
	}

	if level<ZSKIPLIST_MAXLEVEL {
		return level
	}else{
		return ZSKIPLIST_MAXLEVEL
	}
}


/* Insert a new node in the skiplist. Assumes the element does not already
 * exist (up to the caller to enforce that). The skiplist takes ownership
 * of the passed SDS string 'ele'. */
func (zsk *ZSkiplist)ZslInsert(score float64,sds string)*ZSkiplistNode {
	var i,level int
	rank := make([]uint64,ZSKIPLIST_MAXLEVEL)
	update := make([]*ZSkiplistNode,ZSKIPLIST_MAXLEVEL)
	x := zsk.header
	for i= zsk.level-1;i>=0;i++ {
		/* store rank that is crossed to reach the insert position */
		if i== zsk.level-1 {
			rank[i] = 0
		}else{
			rank[i] = rank[i+1]
		}

		for x.levels[i].forward!=nil &&
			(x.levels[i].forward.score<score || (x.levels[i].forward.score==score) && (common.StringBaseCompare(x.levels[i].forward.ele,sds)<0)){
				rank[i] += uint64(x.levels[i].span)
				x = x.levels[i].forward
		}
		update[i] =x
	}

	/* we assume the element is not already inside, since we allow duplicated
	* scores, reinserting the same element should never happen since the
	* caller of zslInsert() should test in the hash table if the element is
	* already inside or not. */

	level = randomLevel()
	if level >zsk.level{
		for k:=zsk.level;k<level;k++ {
			rank[k] = 0
			update[k] = zsk.header
			update[k].levels[k].span = zsk.length
		}
		zsk.level = level
	}

	x = CreateZSkiplistNode(int64(level),score,sds)
	for i := 0;i<level;i++ {
		x.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = x

		/* update span covered by update[i] as x is inserted here */
		x.levels[i].span = update[i].levels[i].span - (rank[0]-rank[i])
		update[i].levels[i].span = (rank[0]-rank[i])+1
	}

	/* increment span for untouched levels */
	for i:= level;i<zsk.level;i++ {
		update[i].levels[i].span ++
	}

	if update[0]==zsk.header{
		x.backward = nil
	}else{
		x.backward = update[0]
	}

	if x.levels[0].forward== nil {
		zsk.tail = x
	}else{
		x.levels[0].forward.backward = x
	}
	zsk.length++
	return x
}

/* Internal function used by zslDelete, zslDeleteByScore and zslDeleteByRank */
func (zsk *ZSkiplist)zslDeleteNode(x *ZSkiplistNode,update []*ZSkiplistNode){
	for i:= 0;i<zsk.level ;i++ {
		if update[i].levels[i].forward == x {
			update[i].levels[i].span  += x.levels[i].span -1
			update[i].levels[i].forward = x.levels[i].forward
		}else{
			update[i].levels[i].span -=1
		}
	}
	if x.levels[0].forward!= nil {
		x.levels[0].forward.backward = x.backward
	}else{
		zsk.tail = x.backward
	}
	for zsk.level> 1 && zsk.header.levels[zsk.level-1].forward == nil {
		zsk.level--
	}
	zsk.length--
}



/* Delete an element with matching score/element from the skiplist.
 * The function returns 1 if the node was found and deleted, otherwise
 * 0 is returned.
 *
 * If 'node' is NULL the deleted node is freed by zslFreeNode(), otherwise
 * it is not freed (but just unlinked) and *node is set to the node pointer,
 * so that it is possible for the caller to reuse the node (including the
 * referenced SDS string at node->ele). */
func (zsk *ZSkiplist)zslDelete(score float64,sds string,node []*ZSkiplistNode)int {
	update := make([]*ZSkiplistNode,ZSKIPLIST_MAXLEVEL)
	x := zsk.header

	for i := zsk.level-1;i >=0;i-- {
		for x.levels[i].forward!=nil &&
			(x.levels[i].forward.score<score || (x.levels[i].forward.score==score && x.levels[i].forward.ele==sds)){

			x = x.levels[i].forward
		}
		update[i] = x
	}

	/* We may have multiple elements with the same score, what we need
     * is to find the element with both the right score and object. */

     x = x.levels[0].forward
     if x!= nil  && score==x.score && x.ele == sds {
     	//zslDeleteNode(zsl, x, update);
     	zsk.zslDeleteNode(x,update)
     	if node==nil {
     		//TODO zslFreeNode(x);
			x=nil
		}else{
			node[0] = x
			return 1

		}

	 }
	/* not found */
	return 0
}

