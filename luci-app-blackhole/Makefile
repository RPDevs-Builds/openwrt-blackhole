include $(TOPDIR)/rules.mk

PKG_NAME:=luci-app-blackhole
PKG_VERSION:=0.1.0
PKG_RELEASE:=1

PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/IamRPDev/luci-app-blackhole.git
PKG_SOURCE_VERSION:=main

include $(INCLUDE_DIR)/package.mk

define Package/luci-app-blackhole
  SECTION:=luci
  CATEGORY:=LuCI
  SUBMENU:=3. Applications
  TITLE:=LuCI support for Blackhole Webserver
  DEPENDS:=+blackhole
  URL:=https://github.com/IamRPDev/luci-app-blackhole
endef

define Package/luci-app-blackhole/description
  Provides a web interface for configuring the Blackhole Webserver.
endef

define Build/Compile
endef

define Package/luci-app-blackhole/install
	$(INSTALL_DIR) $(1)/usr/share/luci/menu.d
	$(INSTALL_DATA) $(PKG_BUILD_DIR)/usr/share/luci/menu.d/luci-app-blackhole.json $(1)/usr/share/luci/menu.d/
	$(INSTALL_DIR) $(1)/usr/share/rpcd/acl.d
	$(INSTALL_DATA) $(PKG_BUILD_DIR)/usr/share/rpcd/acl.d/luci-app-blackhole.json $(1)/usr/share/rpcd/acl.d/
	$(INSTALL_DIR) $(1)/www/luci-static/resources/view
	$(INSTALL_DATA) $(PKG_BUILD_DIR)/www/luci-static/resources/view/blackhole.js $(1)/www/luci-static/resources/view/
endef

$(eval $(call BuildPackage,luci-app-blackhole))
