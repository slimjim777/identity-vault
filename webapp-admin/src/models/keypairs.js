/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */
import Ajax from './Ajax';

var Keypair = {
	url: 'keypairs',

	list: function (account) {
			return Ajax.get(account + '/' + this.url);
	},

	enable:  function(keypairId) {
		return Ajax.post(this.url + '/' + keypairId + '/enable', {});
	},

	disable:  function(keypairId) {
		return Ajax.post(this.url + '/' + keypairId + '/disable', {});
	},

	create:  function(authorityId, key) {
		return Ajax.post(this.url, {'authority-id': authorityId, 'private-key': key});
	},

	generate:  function(authorityId, keyName) {
		return Ajax.post(this.url + '/generate', {'authority-id': authorityId, 'key-name': keyName});
	},

	status:  function() {
		return Ajax.get(this.url + '/status');
	},

	assertion: function(keyId, assertion) {
		return Ajax.post(this.url + '/assertion', {id: keyId, assertion: assertion});
	},

	register: function(keypair) {
		return Ajax.post(this.url + '/register', keypair);
	}
}

export default Keypair;
