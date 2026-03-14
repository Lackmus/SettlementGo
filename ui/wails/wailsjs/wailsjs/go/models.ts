export namespace controllers {
	
	export class CreationOptions {
	    factions: string[];
	    species: string[];
	    traits: string[];
	    npcTypes: string[];
	    npcSubtypeForTypeMap: Record<string, Array<string>>;
	    npcSpeciesForFactionMap: Record<string, Array<string>>;
	
	    static createFrom(source: any = {}) {
	        return new CreationOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.factions = source["factions"];
	        this.species = source["species"];
	        this.traits = source["traits"];
	        this.npcTypes = source["npcTypes"];
	        this.npcSubtypeForTypeMap = source["npcSubtypeForTypeMap"];
	        this.npcSpeciesForFactionMap = source["npcSpeciesForFactionMap"];
	    }
	}

}

export namespace main {
	
	export class SubtypeRoll {
	    stats: string;
	    items: string;
	
	    static createFrom(source: any = {}) {
	        return new SubtypeRoll(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stats = source["stats"];
	        this.items = source["items"];
	    }
	}

}

export namespace mapper {
	
	export class NPCInput {
	    id: string;
	    name: string;
	    type: string;
	    subtype: string;
	    species: string;
	    faction: string;
	    trait: string;
	    stats: string;
	    items: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new NPCInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.subtype = source["subtype"];
	        this.species = source["species"];
	        this.faction = source["faction"];
	        this.trait = source["trait"];
	        this.stats = source["stats"];
	        this.items = source["items"];
	        this.notes = source["notes"];
	    }
	}
	export class SettlementCreateInput {
	    name: string;
	    faction: string;
	    xCoord: number;
	    yCoord: number;
	    population: number;
	    notes: string;
	    initialRandomNpcCount: number;
	
	    static createFrom(source: any = {}) {
	        return new SettlementCreateInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.faction = source["faction"];
	        this.xCoord = source["xCoord"];
	        this.yCoord = source["yCoord"];
	        this.population = source["population"];
	        this.notes = source["notes"];
	        this.initialRandomNpcCount = source["initialRandomNpcCount"];
	    }
	}
	export class SettlementUpdateInput {
	    originalName: string;
	    name: string;
	    faction: string;
	    population: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new SettlementUpdateInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.originalName = source["originalName"];
	        this.name = source["name"];
	        this.faction = source["faction"];
	        this.population = source["population"];
	        this.notes = source["notes"];
	    }
	}
	export class SettlementView {
	    name: string;
	    faction: string;
	    xCoord: number;
	    yCoord: number;
	    population: number;
	    notes: string;
	    npcs: NPCInput[];
	
	    static createFrom(source: any = {}) {
	        return new SettlementView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.faction = source["faction"];
	        this.xCoord = source["xCoord"];
	        this.yCoord = source["yCoord"];
	        this.population = source["population"];
	        this.notes = source["notes"];
	        this.npcs = this.convertValues(source["npcs"], NPCInput);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

