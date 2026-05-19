'use strict';
'require view';
'require form';
'require uci';

return view.extend({
	load: function() {
		return uci.load('blackhole');
	},
	render: function() {
		var m, s, o;

		m = new form.Map('blackhole', _('Blackhole Webserver'), _('Configure the lightweight blackhole server that captures, logs, and mirrors HTTP requests.'));

		s = m.section(form.TypedSection, 'blackhole', _('General Settings'));
		s.anonymous = true;
		s.addremove = false;

		o = s.option(form.Flag, 'enable', _('Enable'));
		o.rmempty = false;

		o = s.option(form.Value, 'port', _('Listen Port'));
		o.default = '8080';
		o.rmempty = false;

		o = s.option(form.Value, 'ip', _('Listen IP Address'));
		o.default = '0.0.0.0';
		o.rmempty = false;

		o = s.option(form.Value, 'log', _('Log File Path'));
		o.default = '/mnt/largedata/blackholeserver/blackhole.log';
		o.rmempty = false;

		o = s.option(form.Value, 'root', _('Mirroring Directory (Root)'));
		o.default = '/mnt/largedata/blackholeserver/mirrored';
		o.rmempty = false;

		o = s.option(form.Value, 'content', _('Content Directory'));
		o.default = '/mnt/largedata/blackholeserver/content';
		o.rmempty = false;

		return m.render();
	}
});
