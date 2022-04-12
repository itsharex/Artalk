import Api from 'artalk/src/api'
import Context from 'artalk/src/context'
import Component from 'artalk/src/lib/component'
import * as Utils from 'artalk/src/lib/utils'
import Comment from 'artalk/src/components/comment'

import SiteList from '../admin/site-list'
import SidebarView from '../sidebar-view'

export default class SitesView extends SidebarView {
  static viewName = 'sites'
  static viewTitle = '站点'
  static viewAdminOnly = true

  viewTabs = {}
  viewActiveTab = ''

  siteList!: SiteList

  mount(siteName: string) {
    if (!this.siteList) {
      this.siteList = new SiteList(this.ctx)
      this.$el.append(this.siteList.$el)
    }

    this.reqSites()
  }

  switchTab(tab: string, siteName: string) {
    this.reqSites()
  }

  async reqSites() {
    const sites = await new Api(this.ctx).siteGet()
    this.siteList.loadSites(sites)
  }
}