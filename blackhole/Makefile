include $(TOPDIR)/rules.mk

PKG_NAME:=blackhole
PKG_VERSION:=0.1.0
PKG_RELEASE:=1

PKG_MAINTAINER:=Gemini CLI <gemini-cli@example.com>
PKG_LICENSE:=MIT
PKG_LICENSE_FILES:=LICENSE

PKG_BUILD_DEPENDS:=golang/host
include $(INCLUDE_DIR)/package.mk

# Link to golang-package.mk
-include $(INCLUDE_DIR)/golang-package.mk

define Package/blackhole
  SECTION:=utils
  CATEGORY:=Utilities
  TITLE:=Blackhole Webserver for Request Tracking
  URL:=https://github.com/IamRPDev/blackhole-server
  DEPENDS:=$(GO_ARCH_DEPENDS)
endef

define Package/blackhole/description
  A lightweight Go webserver that logs all requests and mirrors their paths to the filesystem.
  Returns a 1x1 transparent tracking GIF for all requests.
endef

GO_PKG:=github.com/IamRPDev/blackhole-server
GO_PKG_BUILD_PKG:=$(GO_PKG)/cmd/blackhole

define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
	$(CP) ./src/* $(PKG_BUILD_DIR)/
	$(CP) ./files $(PKG_BUILD_DIR)/
endef

# Explicitly define compilation to avoid any ambiguity with golang-package.mk versions
define Build/Compile
	cd $(PKG_BUILD_DIR) && \
	GOOS=linux \
	GOARCH=$(ARCH) \
	$(STAGING_DIR_HOST)/bin/go build -v -o blackhole ./cmd/blackhole
endef

define Package/blackhole/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/blackhole $(1)/usr/bin/blackhole
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) $(PKG_BUILD_DIR)/files/blackhole.config $(1)/etc/config/blackhole
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/files/blackhole.init $(1)/etc/init.d/blackhole
endef

$(eval $(call BuildPackage,blackhole))
