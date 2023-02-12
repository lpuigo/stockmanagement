package movementeditmodal

const template string = `<el-dialog
        :before-close="HideWithControl"
        :visible.sync="Visible" :close-on-click-modal="false"
        width="70vw" top="10vh"
>
    <!-- 
        Modal Title
    -->
    <div slot="title">
		<h2 style="margin: 0 10px" v-if="current_movement">
			<i class="far fa-edit icon--left"></i>Edition du Hébergement : <span style="color: #ccebff">{{edited_movement.Ref}} {{FormatAddress(edited_movement)}}</span>
		</h2>
    </div>
    <!-- 
        Modal Body
        style="height: 100%;"
    -->
	<el-tabs tab-position="top" v-model="EditMode" style="height: 50vh;margin: 0px 15px;">
		<!-- ====================	Movement Edition  ============================== -->
		<el-tab-pane :label="current_movement.Type" name="acc" lazy=true style="padding: 0px 0px;">
			< style="margin: 0px 20px;">
				
				<!--	Actor & Responsible -->
				<el-row :gutter="5" type="flex" align="middle" class="spaced">
					<el-col :span="2" class="align-right">Acteur:</el-col>
					<el-col :span="10">
						<el-input v-model="current_movement.Actor" size="mini""></el-input>
					</el-col>
					<el-col :span="2" class="align-right">Responsable:</el-col>
					<el-col :span="10">
						<el-input v-model="current_movement.Responsible" size="mini""></el-input>
					</el-col>
				</el-row>
		
				<!--	Date  -->
				<el-row :gutter="5" type="flex" align="middle" class="doublespaced">
					<el-col :span="2" class="align-right">Date:</el-col>
					<el-col :span="10">
	                    <el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
							style="width: 100%" type="date"
							v-model="current_movement.Date"
							value-format="yyyy-MM-dd"
							:disabled="DisableUpdateDate"
							:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
							@change="UpdateDate"
    	                ></el-date-picker>
					</el-col>
				</el-row>
				
				<!--	Article Flows -->
				<el-row :gutter="5" type="flex" align="middle" class="doublespaced">
					<el-col :span="2" class="align-right">Liste d'articles:</el-col>
				</el-row>

		
				<!--	Address & Zipcode -->
				<el-row :gutter="5" type="flex" align="middle" class="spaced">
					<el-col :span="2" class="align-right">Adresse:</el-col>
					<el-col :span="10">
						<el-input v-model="current_movement.Address" size="mini" @change="CheckInputCase(current_movement)"></el-input>
					</el-col>
					<el-col :span="2" class="align-right">Code Postal:</el-col>
					<el-col :span="9">
						<el-input v-model="current_movement.Zipcode" size="mini"></el-input>
					</el-col>
					<el-col :span="1">
						<el-tooltip content="Copier l'adresse dans le presse-papier" placement="bottom" effect="light" open-delay="500">
							<el-button type="info" plain style="width: 100%" class="icon" icon="fa-regular fa-clipboard" size="mini" @click="CopyAddressToClipBoard(current_movement)"></el-button>
						</el-tooltip>
					</el-col>
				</el-row>
		
				<!--	Lat & Long -->		
				<el-row :gutter="5" type="flex" align="middle" class="doublespaced">
					<el-col :span="2" class="align-right">Lat, Long:</el-col>
					<el-col :span="4">
						<el-input v-model="current_latlong" size="mini" @input="UpdateLatLong"></el-input>
					</el-col>
					<el-col :span="1">
						<el-tooltip content="Montrer la position sur Google Maps" placement="bottom" effect="light" open-delay="500">
							<el-button type="info" plain style="width: 100%" class="icon" icon="fa-solid fa-location-dot" size="mini" @click="GMaps(current_latlong)"></el-button>
						</el-tooltip>
					</el-col>
					<el-col :span="1">
						<el-tooltip content="Rechercher la position d'après l'adresse" placement="bottom" effect="light" open-delay="500">
							<el-button type="info" plain style="width: 100%" class="icon" icon="fa-solid fa-magnifying-glass-location" size="mini" @click="SearchByAddress(current_movement)"></el-button>
						</el-tooltip>
					</el-col>
				</el-row>
		
				<!--	Owner Name & Phone -->
				<el-row :gutter="5" type="flex" align="middle" class="doublespaced">
					<el-col :span="2" class="align-right">Propriétaire:</el-col>
					<el-col :span="10">
						<el-input v-model="current_movement.OwnerName" size="mini"></el-input>
					</el-col>
					<el-col :span="2" class="align-right">Téléphone:</el-col>
					<el-col :span="10">
						<el-input v-model="current_movement.OwnerPhone" size="mini"></el-input>
					</el-col>
				</el-row>
					
				<!--	Owner Mail Address -->
				<el-row :gutter="5" type="flex" align="middle" class="spaced">
					<el-col :offset="12" :span="2" class="align-right">@Mail:</el-col>
					<el-col :span="10">
						<el-input v-model="current_movement.OwnerMail" size="mini" @change="CheckInputCase(current_movement)"></el-input>
					</el-col>
				</el-row>
		
		
			</div>
		</el-tab-pane>

		<!-- ====================	Rental Stays Edition  ============================== -->
		<el-tab-pane label="Séjours" name="stay" lazy=true style="padding: 0px 20px;height: 40vh">
			<pre>{{current_movement}}</pre>
		</el-tab-pane>

	</el-tabs>

    <!-- 
        Modal Footer Action Bar
    -->
    <span slot="footer">
		<el-row :gutter="15">
			<el-col :span="24" style="text-align: right">
				<el-tooltip :open-delay="500" effect="light">
					<div slot="content">Annuler les changements</div>
					<el-button :disabled="!hasChanged" @click="UndoChange" icon="fas fa-undo-alt" plain size="mini"
                               type="info"></el-button>
				</el-tooltip>
				
				<el-button @click="Hide" size="mini">Fermer</el-button>
				
				<el-button :disabled="!hasChanged" type="success" @click="ConfirmChange" plain size="mini">Valider</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`
