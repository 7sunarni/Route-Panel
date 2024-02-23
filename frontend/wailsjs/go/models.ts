export namespace route {
	
	export class Route {
	    destination: string;
	    mask: string;
	    gateway: string;
	    interfaceIp: string;
	    interfaceName: string;
	    metric: string;
	    type: string;
	    protocol: string;
	
	    static createFrom(source: any = {}) {
	        return new Route(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.destination = source["destination"];
	        this.mask = source["mask"];
	        this.gateway = source["gateway"];
	        this.interfaceIp = source["interfaceIp"];
	        this.interfaceName = source["interfaceName"];
	        this.metric = source["metric"];
	        this.type = source["type"];
	        this.protocol = source["protocol"];
	    }
	}

}

