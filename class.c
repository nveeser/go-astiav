#include <libavutil/log.h>
#include <stdint.h>
#include <stdlib.h>

char* astiavClassItemName(AVClass* c, void* ptr) {
	return (char*)c->item_name(ptr);
}

AVClassCategory astiavClassCategory(AVClass* c, void* ptr) {
	if (c->get_category) return c->get_category(ptr);
	return c->category;
}

AVClass** astiavClassParent(AVClass* c, void* ptr) {
	if (c->parent_log_context_offset) {
	    uint8_t* pptr = ((uint8_t *) ptr);
		AVClass** parent = *(AVClass ***) (pptr + c->parent_log_context_offset);
		if (parent && *parent) {
			return parent;
		}
	}
	return NULL;
}