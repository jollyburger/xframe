#=============================================================================
#     FileName: Makefile
#         Desc: api server makefile
#       Author: ato.ye
#        Email: ato.ye@ucloud.cn
#     HomePage: http://www.ucloud.cn
#      Version: 0.0.1
#   LastChange: 2016-04-20 20:13:59
#      History:
#=============================================================================

SRCDIRS_EXCLUDE = proto log
SRCDIRS_ALL = $(sort $(subst /,,$(dir $(wildcard */))))
SRCDIRS = $(filter-out $(SRCDIRS_EXCLUDE), $(SRCDIRS_ALL))

PKGDIRS_EXCLUDE=$(GOROOT)/pkg
PKGDIRS_ALL = $(addsuffix /pkg, $(subst :, ,$(GOPATH)))
PKGDIRS = $(filter-out $(PKGDIRS_EXCLUDE), $(PKGDIRS_ALL))

all: build_main 
	@for subdir in $(SRCDIRS);do \
		cd $$subdir; go install; cd ..; \
	done 

build_main:
	@go build echo_server.go;

show:
	@echo "==================src====================="
	@echo SRCDIRS_ALL: $(SRCDIRS_ALL)
	@echo SRCDIRS_EXCLUDE: $(SRCDIRS_EXCLUDE)
	@echo SRCDIRS: $(SRCDIRS)
	@echo "==================pkg====================="
	@echo PKGDIRS_EXCLUDE: $(PKGDIRS_EXCLUDE)
	@echo PKGDIRS_ALL: $(PKGDIRS_ALL)
	@echo PKGDIRS: $(PKGDIRS)
	@echo "================clean====================="
	@for subdir in $(PKGDIRS); do \
		cd $$subdir;\
		module_name=`echo $(CURDIR)|awk -F"/" '{print $$(NF)}'`;\
		result=`find . |grep $$module_name |head -n1|awk -F"." '{print $$2}'`; \
		if [ -n "$$result" ];then \
			echo clean_dirs:$$subdir$$result; \
		fi \
	done

clean:
	@rm -f echo_server
	@for subdir in $(PKGDIRS); do \
		cd $$subdir;\
		module_name=`echo $(CURDIR)|awk -F"/" '{print $$(NF)}'`;\
		result=`find . |grep $$module_name|head -n1|awk -F"." '{print $$2}'`; \
		if [ -n "$$result" ];then \
			cd $$subdir$$result; rm -rf *; \
		fi \
	done
