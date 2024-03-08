export namespace route {
	
	export class Route {
	    destination: string;
	    mask: string;
	    gateway: string;
	    interfaceIp: string;
	    interfaceName: string;
	
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
	    }
	}

}

